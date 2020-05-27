package channel

import (
	"github.com/yanlong-li/HelloWorld-GO/io/logger"
	"github.com/yanlong-li/HelloWorld-GO/io/network/connect"
	"github.com/yanlong-li/HelloWorld-GO/io/network/route"
	"github.com/yanlong-li/HelloWorldServer/model"
	"github.com/yanlong-li/HelloWorldServer/model/online"
	"github.com/yanlong-li/HelloWorldServer/packetModel/server/channel"
)

func init() {
	route.Register(channel.GetList{}, actionGetChannelList)
}

func actionGetChannelList(_ channel.GetList, conn connect.Connector) {

	user, _ := online.Auth(conn.GetId())

	// model
	_list := model.GetUserChannels(user.Id)
	// packet
	list := channel.List{
		List: make([]channel.Info, 0),
	}
	for _, cu := range _list {

		_cu, ok := cu.(model.ChannelUser)

		if ok != true {
			logger.Fatal("断言错误", 0, cu)
			continue
		}

		cha, err := _cu.Channel()
		if !err.Status() {
			continue
		}
		createUser, _ := model.GetUserById(cha.CreateUserId)
		ownerUser, _ := model.GetUserById(cha.CreateUserId)
		info := channel.Info{
			Id:     cha.Id,
			Name:   cha.Name,
			Verify: true,
			CreateUser: struct {
				Id       uint64
				Nickname string
			}{Id: createUser.Id, Nickname: createUser.Nickname},
			OwnerUser: struct {
				Id       uint64
				Nickname string
			}{Id: ownerUser.Id, Nickname: ownerUser.Nickname},
			CreateTime: cha.CreateTime,
			Public:     true,
			Avatar:     cha.Avatar,
			Describe:   cha.Describe,
			Channels:   []channel.Info{},
		}
		ChannelChildrenS := cha.GetChildren()
		for _, channelChildren := range ChannelChildrenS {
			if _channelChildren, ok := channelChildren.(model.Channel); ok {
				_createUser, _ := model.GetUserById(cha.CreateUserId)
				_ownerUser, _ := model.GetUserById(cha.CreateUserId)
				_channelInfo := channel.Info{
					Id:     _channelChildren.Id,
					Name:   _channelChildren.Name,
					Verify: true,
					CreateUser: struct {
						Id       uint64
						Nickname string
					}{Id: _createUser.Id, Nickname: _createUser.Nickname},
					OwnerUser: struct {
						Id       uint64
						Nickname string
					}{Id: _ownerUser.Id, Nickname: _ownerUser.Nickname},
					CreateTime: _channelChildren.CreateTime,
					Public:     true,
					Avatar:     _channelChildren.Avatar,
					Describe:   _channelChildren.Describe,
					Channels:   []channel.Info{},
				}
				info.Channels = append(info.Channels, _channelInfo)
			}

		}

		list.List = append(list.List, info)
	}

	_ = conn.Send(list)
}
