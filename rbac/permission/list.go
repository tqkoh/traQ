package permission

import "github.com/mikespook/gorbac"

// 全パーミッションのリスト。パーミッションを新たに定義した場合はここに必ず追加すること
var list = map[string]gorbac.Permission{
	CreateChannel.ID(): CreateChannel,
	GetChannel.ID():    GetChannel,
	EditChannel.ID():   EditChannel,
	DeleteChannel.ID(): DeleteChannel,

	GetTopic.ID():  GetTopic,
	EditTopic.ID(): EditTopic,

	GetMessage.ID():    GetMessage,
	PostMessage.ID():   PostMessage,
	EditMessage.ID():   EditMessage,
	DeleteMessage.ID(): DeleteMessage,

	GetPin.ID():    GetPin,
	CreatePin.ID(): CreatePin,
	DeletePin.ID(): DeletePin,

	GetNotificationStatus.ID():     GetNotificationStatus,
	ChangeNotificationStatus.ID():  ChangeNotificationStatus,
	ConnectNotificationStream.ID(): ConnectNotificationStream,
	RegisterDevice.ID():            RegisterDevice,

	GetUser.ID():      GetUser,
	GetMe.ID():        GetMe,
	RegisterUser.ID(): RegisterUser,
	EditMe.ID():       EditMe,
	ChangeMyIcon.ID(): ChangeMyIcon,

	GetClip.ID():    GetClip,
	CreateClip.ID(): CreateClip,
	DeleteClip.ID(): DeleteClip,

	GetStar.ID():    GetStar,
	CreateStar.ID(): CreateStar,
	DeleteStar.ID(): DeleteStar,

	GetChannelVisibility.ID():    GetChannelVisibility,
	ChangeChannelVisibility.ID(): ChangeChannelVisibility,

	GetUnread.ID():    GetUnread,
	DeleteUnread.ID(): DeleteUnread,

	GetTag.ID():             GetTag,
	AddTag.ID():             AddTag,
	RemoveTag.ID():          RemoveTag,
	ChangeTagLockState.ID(): ChangeTagLockState,

	GetStamp.ID():           GetStamp,
	CreateStamp.ID():        CreateStamp,
	EditStamp.ID():          EditStamp,
	DeleteStamp.ID():        DeleteStamp,
	GetMessageStamp.ID():    GetMessageStamp,
	AddMessageStamp.ID():    AddMessageStamp,
	RemoveMessageStamp.ID(): RemoveMessageStamp,

	UploadFile.ID():   UploadFile,
	DownloadFile.ID(): DownloadFile,
	DeleteFile.ID():   DeleteFile,

	GetHeartbeat.ID():  GetHeartbeat,
	PostHeartbeat.ID(): PostHeartbeat,

	GetWebhook.ID():    GetWebhook,
	CreateWebhook.ID(): CreateWebhook,
	EditWebhook.ID():   EditWebhook,
	DeleteWebhook.ID(): DeleteWebhook,
}

// GetPermission : パーミッション名からgorbac.Permissionを取得します
func GetPermission(name string) gorbac.Permission {
	return list[name]
}

// GetAllPermissionList : 全パーミッションリストを返します
func GetAllPermissionList() map[string]gorbac.Permission {
	return list
}
