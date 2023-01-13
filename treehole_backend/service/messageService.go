package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
	"treehole/dao"
	"treehole/define"
	"treehole/models"
	"treehole/utils"

	"treehole/logger"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

/**
 * @Author jiang
 * @Description 聊天室
 * @Date 13:00 2023/1/13
 **/
type Node struct {
	Conn          *websocket.Conn // 连接
	HeartbeatTime uint64          // 心跳时间
}

//更新用户心跳
func (node *Node) Heartbeat(currentTime uint64) {
	node.HeartbeatTime = currentTime
}

// 映射关系
var clientMap map[string]*Node = make(map[string]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

// 防止跨域站点伪造请求
// 将HTTP升级成Websocket的全局变量upgrade，并默认允许跨域
var upGrade = websocket.Upgrader{
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebsocketMessage(c *gin.Context) {
	// 1、升级协议为websocket
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.SugarLogger.Error("Websocket Upgrade Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "系统异常")
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			logger.SugarLogger.Error("Websocket Colse Error:" + err.Error())
		}
	}(ws)

	// 2、读取token上的用户信息，将数据放入缓存
	u, _ := c.Get("user_claims")
	userClaim := u.(*utils.UserClaims)

	user := models.UserBasic{
		Username: userClaim.Username,
		Identity: userClaim.Identity,
		Usericon: userClaim.Usericon,
	}
	err = dao.RedisMessageOnlineCreate(c, user)
	if err != nil {
		logger.SugarLogger.Error("RedisMessageOnlineCreate Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "进入聊天教室失败")
		return
	}

	// 3、创建node
	currentTime := uint64(time.Now().Unix())
	node := &Node{
		Conn:          ws,
		HeartbeatTime: currentTime,
	}
	rwLocker.Lock()
	clientMap[userClaim.Identity] = node
	rwLocker.Unlock()

	// 4、发送、接收信息
	for {
		// 1）接收消息
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			logger.SugarLogger.Error("Read Error:" + err.Error())
			return
		}

		// 2）将消息转换为json格式，存入缓存、数据库
		msg := models.MessageStruct{}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			logger.SugarLogger.Error("Read json.Unmarshal Error:" + err.Error())
			return
		}

		// 3）判断是心跳检测类型还是消息类型
		if msg.Type == -1 {
			if _, ok := clientMap[userClaim.Identity]; ok {
				node.Conn.WriteMessage(websocket.TextMessage, []byte(data))
			}
		} else {
			// 心跳检测信息，更新心跳时间
			currentTime := uint64(time.Now().Unix())
			node.Heartbeat(currentTime)

			// 4）消息类型，先将聊天记录存入缓存和数据库当中，后发送
			// 判断identity是否存在
			cnt, err := dao.FindUserByIdentityCount(msg.UserIdentity)
			if err != nil {
				logger.SugarLogger.Error("FindUserByIdentityCount Error:" + err.Error())
				utils.RespFail(c, int(define.FailCode), "发送消息失败")
				return
			}

			if cnt <= 0 {
				logger.SugarLogger.Error("FindUserByIdentityCount Error: ")
				utils.RespFail(c, int(define.FailCode), "发送消息失败")
				return
			}

			// 判断useridentity是否为自己，如果是，则表示为自己发送，需要存入数据库
			if userClaim.Identity == msg.UserIdentity {
				identity, _ := utils.GetID()
				messageBasic := models.MessageBasic{
					Identity:     strconv.Itoa(int(identity)),
					UserIdentity: msg.UserIdentity,
					Username:     msg.Username,
					Usericon:     msg.Usericon,
					Date:         msg.Date,
					CreateTime:   msg.CreateTime,
				}
				// 将数据存入数据库
				err = dao.CreateMessageBasic(messageBasic)
				if err != nil {
					logger.SugarLogger.Error("CreateMessageBasic Error:" + err.Error())
					utils.RespFail(c, int(define.FailCode), "发送消息失败")
					return
				}

				// 将数据存入缓存
				err = dao.RedisCreateMessage(c, messageBasic)
				if err != nil {
					logger.SugarLogger.Error("RedisCreateMessage Error:" + err.Error())
					utils.RespFail(c, int(define.FailCode), "发送消息失败")
					return
				}
			}

			// 发送消息
			for _, v := range clientMap {
				v.Conn.WriteMessage(websocket.TextMessage, []byte(data))
				if err != nil {
					logger.SugarLogger.Error("WriteMessage Error:" + err.Error())
					return
				}
			}
		}
	}
}

// 清理超时连接
func CleanConnection(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("cleanConnection err", r)
		}
	}()

	currentTime := uint64(time.Now().Unix())
	for i := range clientMap {
		node := clientMap[i]
		if node.IsHeartbeatTimeOut(currentTime) {
			fmt.Println("心跳超时..... 关闭连接：", node)
			logger.SugarLogger.Info("心跳超时..... 关闭连接：", node)
			node.Conn.Close()
		}
	}
	return result
}

// 判断用户心跳是否超时
func (node *Node) IsHeartbeatTimeOut(currentTime uint64) (timeout bool) {
	if node.HeartbeatTime+viper.GetUint64("timeout.HeartbeatMaxTime") <= currentTime {
		fmt.Println("心跳超时。。。自动下线", node)
		logger.SugarLogger.Info("心跳超时。。。自动下线", node)
		timeout = true
	}
	return
}

// GetOnlistList
// @Summary 获取在线人数
// @Tags 聊天室接口
// @Param Authorization header string true "Authorization"
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {object} utils.H
// @Router /chat/websocket/getOnlineList [get]
func GetOnlistList(c *gin.Context) {
	// 1、获取参数
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	page = (page - 1) * size
	count := 0

	// 2、获取用户信息的在线list信息
	ListKey := viper.GetString("redis.KeyWebsocketOnlineList")
	userList := utils.RDB.LRange(c, ListKey, 0, -1).Val()

	list := []map[string]string{}
	for i, userid := range userList {
		// 我们需要的分页数据
		if i >= page && i < page+size {
			// 3、判断缓存中用户信息是否已经过期，不过期则查询
			HashKey := viper.GetString("redis.KeyWebsocketOnlineHashPrefix") + userid

			if utils.RDB.Exists(c, HashKey).Val() == 0 {
				// 4、过期，则到数据库提取用户信息，同时将期存入缓存，过期时间为1天
				// 从数据库中提取用户信息
				user, err := dao.FindByIdentity(userid)
				if err != nil {
					logger.SugarLogger.Error("FindByIdentity Error:" + err.Error())
					utils.RespFail(c, int(define.FailCode), "获取在线人数失败")
					return
				}

				// 插入缓存
				userInfo := map[string]string{
					"useridentity": user.Identity,
					"usericon":     user.Usericon,
					"username":     user.Username,
				}
				list = append(list, userInfo)
				HashKey := viper.GetString("redis.KeyWebsocketMessageHashPrefix") + user.Identity
				fmt.Println("HashKey：", HashKey)

				err = dao.RedisSave(c, HashKey, userInfo)
				if err != nil {
					logger.SugarLogger.Error("RedisSave Error:" + err.Error())
					utils.RespFail(c, int(define.FailCode), "获取在线人数失败")
					return
				}
			} else {
				info := utils.RDB.HGetAll(c, HashKey).Val()
				list = append(list, info)
			}
		}
		count++
	}

	data := map[string]interface{}{
		"count":      count,
		"onlineInfo": list,
	}

	utils.RespSuccess(c, "获取在线人数信息成功", data)
}

// GetMessageList
// @Summary 获取聊天记录
// @Tags 聊天室接口
// @Param Authorization header string true "Authorization"
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {object} utils.H
// @Router /chat/websocket/getMessageList [get]
func GetMessageList(c *gin.Context) {
	// 1、获取参数
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	page = (page - 1) * size

	messageList := []map[string]string{}
	// 	1、判断查询的查询的大小，判断查询大小是否在500以内。
	if page+size >= 500 {
		// 2、500以内获取缓存数据，返回。
		ListKey := viper.GetString("redis.KeyWebsocketMessageList")

		identityList := utils.RDB.LRange(c, ListKey, int64(page), int64(page+size-1)).Val()

		for _, identity := range identityList {
			HashKey := viper.GetString("redis.KeyWebsocketMessageHashPrefix") + identity
			list := utils.RDB.HGetAll(c, HashKey).Val()
			messageList = append(messageList, list)
		}
	} else {
		// 3、500以外查询数据库，返回。
		tx := dao.FindMessageBasic()
		List := []models.MessageBasic{}
		err := tx.Offset(page).Limit(size).Find(&List).Error
		if err != nil {
			logger.SugarLogger.Error("FindMessage Error:" + err.Error())
			utils.RespFail(c, int(define.FailCode), "获取聊天记录失败")
			return
		}

		for _, message := range List {
			messageInfo := map[string]string{
				"useridentity": message.UserIdentity,
				"usericon":     message.Usericon,
				"username":     message.Username,
				"date":         message.Date,
				"create_time":  message.CreateTime,
			}
			messageList = append(messageList, messageInfo)
		}

	}

	data := map[string]interface{}{
		"messageList": messageList,
	}

	utils.RespSuccess(c, "查询聊天记录成功", data)
}
