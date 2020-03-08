package v3

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traQ/model"
	"time"
)

type Channel struct {
	ID         uuid.UUID     `json:"id"`
	Name       string        `json:"name"`
	ParentID   uuid.NullUUID `json:"parentId"`
	Topic      string        `json:"topic"`
	Children   []uuid.UUID   `json:"children"`
	Visibility bool          `json:"visibility"`
	Force      bool          `json:"force"`
}

func formatChannel(channel *model.Channel, childrenID []uuid.UUID) *Channel {
	return &Channel{
		ID:         channel.ID,
		Name:       channel.Name,
		ParentID:   uuid.NullUUID{UUID: channel.ParentID, Valid: channel.ParentID != uuid.Nil},
		Topic:      channel.Topic,
		Children:   childrenID,
		Visibility: channel.IsVisible,
		Force:      channel.IsForced,
	}
}

type UserTag struct {
	ID        uuid.UUID `json:"tagId"`
	Tag       string    `json:"tag"`
	IsLocked  bool      `json:"isLocked"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func formatUserTags(uts []*model.UsersTag) []UserTag {
	res := make([]UserTag, len(uts))
	for i, ut := range uts {
		res[i] = UserTag{
			ID:        ut.Tag.ID,
			Tag:       ut.Tag.Name,
			IsLocked:  ut.IsLocked,
			CreatedAt: ut.CreatedAt,
			UpdatedAt: ut.UpdatedAt,
		}
	}
	return res
}

type User struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName"`
	IconFileID  uuid.UUID `json:"iconFileId"`
	Bot         bool      `json:"bot"`
	State       int       `json:"state"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func formatUsers(users []*model.User) []User {
	res := make([]User, len(users))
	for i, user := range users {
		res[i] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.GetResponseDisplayName(),
			IconFileID:  user.Icon,
			Bot:         user.Bot,
			State:       user.Status.Int(),
			UpdatedAt:   user.UpdatedAt,
		}
	}
	return res
}

type UserDetail struct {
	ID          uuid.UUID   `json:"id"`
	State       int         `json:"state"`
	Bot         bool        `json:"bot"`
	IconFileID  uuid.UUID   `json:"iconFileId"`
	DisplayName string      `json:"displayName"`
	Name        string      `json:"name"`
	TwitterID   string      `json:"twitterId"`
	LastOnline  *time.Time  `json:"lastOnline"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	Tags        []UserTag   `json:"tags"`
	Groups      []uuid.UUID `json:"groups"`
	Bio         string      `json:"bio"`
}

func formatUserDetail(user *model.User, uts []*model.UsersTag, g []uuid.UUID) *UserDetail {
	return &UserDetail{
		ID:          user.ID,
		State:       user.Status.Int(),
		Bot:         user.Bot,
		IconFileID:  user.Icon,
		DisplayName: user.GetResponseDisplayName(),
		Name:        user.Name,
		TwitterID:   user.TwitterID,
		LastOnline:  user.LastOnline.Ptr(),
		UpdatedAt:   user.UpdatedAt,
		Tags:        formatUserTags(uts),
		Groups:      g,
		Bio:         "", // TODO
	}
}

type Webhook struct {
	WebhookID   string    `json:"id"`
	BotUserID   string    `json:"botUserId"`
	DisplayName string    `json:"displayName"`
	Description string    `json:"description"`
	Secure      bool      `json:"secure"`
	ChannelID   string    `json:"channelId"`
	OwnerID     string    `json:"ownerId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func formatWebhook(w model.Webhook) *Webhook {
	return &Webhook{
		WebhookID:   w.GetID().String(),
		BotUserID:   w.GetBotUserID().String(),
		DisplayName: w.GetName(),
		Description: w.GetDescription(),
		Secure:      len(w.GetSecret()) > 0,
		ChannelID:   w.GetChannelID().String(),
		OwnerID:     w.GetCreatorID().String(),
		CreatedAt:   w.GetCreatedAt(),
		UpdatedAt:   w.GetUpdatedAt(),
	}
}

func formatWebhooks(ws []model.Webhook) []*Webhook {
	res := make([]*Webhook, len(ws))
	for i, w := range ws {
		res[i] = formatWebhook(w)
	}
	return res
}

type Bot struct {
	ID              uuid.UUID       `json:"id"`
	BotUserID       uuid.UUID       `json:"botUserId"`
	Description     string          `json:"description"`
	DeveloperID     uuid.UUID       `json:"developerId"`
	SubscribeEvents model.BotEvents `json:"subscribeEvents"`
	State           model.BotState  `json:"state"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}

func formatBot(b *model.Bot) *Bot {
	return &Bot{
		ID:              b.ID,
		BotUserID:       b.BotUserID,
		Description:     b.Description,
		SubscribeEvents: b.SubscribeEvents,
		State:           b.State,
		DeveloperID:     b.CreatorID,
		CreatedAt:       b.CreatedAt,
		UpdatedAt:       b.UpdatedAt,
	}
}

func formatBots(bs []*model.Bot) []*Bot {
	res := make([]*Bot, len(bs))
	for i, b := range bs {
		res[i] = formatBot(b)
	}
	return res
}

type BotTokens struct {
	VerificationToken string `json:"verificationToken"`
	AccessToken       string `json:"accessToken"`
}

type BotDetail struct {
	ID              uuid.UUID       `json:"id"`
	BotUserID       uuid.UUID       `json:"botUserId"`
	Description     string          `json:"description"`
	DeveloperID     uuid.UUID       `json:"developerId"`
	SubscribeEvents model.BotEvents `json:"subscribeEvents"`
	State           model.BotState  `json:"state"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
	Tokens          BotTokens       `json:"tokens"`
	Endpoint        string          `json:"endpoint"`
	Privileged      bool            `json:"privileged"`
	Channels        []uuid.UUID     `json:"channels"`
}

func formatBotDetail(b *model.Bot, t *model.OAuth2Token, channels []uuid.UUID) *BotDetail {
	return &BotDetail{
		ID:              b.ID,
		BotUserID:       b.BotUserID,
		Description:     b.Description,
		SubscribeEvents: b.SubscribeEvents,
		State:           b.State,
		DeveloperID:     b.CreatorID,
		CreatedAt:       b.CreatedAt,
		UpdatedAt:       b.UpdatedAt,
		Tokens: BotTokens{
			VerificationToken: b.VerificationToken,
			AccessToken:       t.AccessToken,
		},
		Endpoint:   b.PostURL,
		Privileged: b.Privileged,
		Channels:   channels,
	}
}

type Message struct {
	ID        uuid.UUID            `json:"id"`
	UserID    uuid.UUID            `json:"userId"`
	ChannelID uuid.UUID            `json:"channelId"`
	Content   string               `json:"content"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
	Pinned    bool                 `json:"pinned"`
	Stamps    []model.MessageStamp `json:"stamps"`
	ThreadID  uuid.NullUUID        `json:"threadId"` // TODO
}

func formatMessage(m *model.Message) *Message {
	return &Message{
		ID:        m.ID,
		UserID:    m.UserID,
		ChannelID: m.ChannelID,
		Content:   m.Text,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Pinned:    m.Pin != nil,
		Stamps:    m.Stamps,
	}
}

func formatMessages(ms []*model.Message) []*Message {
	res := make([]*Message, len(ms))
	for i, m := range ms {
		res[i] = formatMessage(m)
	}
	return res
}

type Pin struct {
	UserID   uuid.UUID `json:"userId"`
	PinnedAt time.Time `json:"pinnedAt"`
	Message  *Message  `json:"message"`
}

func formatPin(pin *model.Pin) *Pin {
	res := &Pin{
		UserID:   pin.UserID,
		PinnedAt: pin.CreatedAt,
		Message:  formatMessage(&pin.Message),
	}
	res.Message.Pinned = true
	return res
}

func formatPins(pins []*model.Pin) []*Pin {
	res := make([]*Pin, len(pins))
	for i, p := range pins {
		res[i] = formatPin(p)
	}
	return res
}

type MessagePin struct {
	UserID   uuid.UUID `json:"userId"`
	PinnedAt time.Time `json:"pinnedAt"`
}

func formatMessagePin(pin *model.Pin) *Pin {
	return &Pin{
		UserID:   pin.UserID,
		PinnedAt: pin.CreatedAt,
	}
}

type UserGroupMember struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"role"`
}

func formatUserGroupMembers(members []*model.UserGroupMember) []UserGroupMember {
	arr := make([]UserGroupMember, len(members))
	for i, m := range members {
		arr[i] = UserGroupMember{
			ID:   m.UserID,
			Role: m.Role,
		}
	}
	return arr
}

func formatUserGroupAdmins(admins []*model.UserGroupAdmin) []uuid.UUID {
	arr := make([]uuid.UUID, len(admins))
	for i, m := range admins {
		arr[i] = m.UserID
	}
	return arr
}

type UserGroup struct {
	ID          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        string            `json:"type"`
	Members     []UserGroupMember `json:"members"`
	Admins      []uuid.UUID       `json:"admins"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

func formatUserGroup(g *model.UserGroup) *UserGroup {
	return &UserGroup{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		Type:        g.Type,
		Members:     formatUserGroupMembers(g.Members),
		Admins:      formatUserGroupAdmins(g.Admins),
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
}

func formatUserGroups(gs []*model.UserGroup) []*UserGroup {
	arr := make([]*UserGroup, len(gs))
	for i, g := range gs {
		arr[i] = formatUserGroup(g)
	}
	return arr
}
