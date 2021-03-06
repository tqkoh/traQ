package repository

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/traQ/model"
	"github.com/traPtitech/traQ/utils/optional"
	"testing"
)

func TestGormRepository_UpdateChannel(t *testing.T) {
	t.Parallel()
	repo, _, _, user := setupWithUser(t, common)

	cases := []UpdateChannelArgs{
		{
			UpdaterID: user.GetID(),
			Topic:     optional.StringFrom("test"),
		},
		{
			UpdaterID: user.GetID(),
			Topic:     optional.StringFrom(""),
		},
		{
			UpdaterID:          user.GetID(),
			Visibility:         optional.BoolFrom(true),
			ForcedNotification: optional.BoolFrom(true),
		},
		{
			UpdaterID:          user.GetID(),
			Visibility:         optional.BoolFrom(true),
			ForcedNotification: optional.BoolFrom(false),
		},
		{
			UpdaterID:          user.GetID(),
			Visibility:         optional.BoolFrom(false),
			ForcedNotification: optional.BoolFrom(true),
		},
		{
			UpdaterID:          user.GetID(),
			Visibility:         optional.BoolFrom(false),
			ForcedNotification: optional.BoolFrom(false),
		},
	}

	for i, v := range cases {
		v := v
		i := i
		t.Run(fmt.Sprintf("Case%d", i), func(t *testing.T) {
			t.Parallel()
			ch := mustMakeChannel(t, repo, rand)
			changed, err := repo.UpdateChannel(ch.ID, v)
			if assert.NoError(t, err) {
				ch, err := repo.GetChannel(ch.ID)
				require.NoError(t, err)
				assert.EqualValues(t, ch, changed)
			}
		})
	}
}

func TestRepositoryImpl_GetChannel(t *testing.T) {
	t.Parallel()
	repo, _, _ := setup(t, common)
	channel := mustMakeChannel(t, repo, rand)

	t.Run("Exists", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		ch, err := repo.GetChannel(channel.ID)
		if assert.NoError(err) {
			assert.Equal(channel.ID, ch.ID)
			assert.Equal(channel.Name, ch.Name)
		}
	})

	t.Run("NotExists", func(t *testing.T) {
		_, err := repo.GetChannel(uuid.Nil)
		assert.Error(t, err)
	})
}

func TestGormRepository_ChangeChannelSubscription(t *testing.T) {
	t.Parallel()
	repo, _, _ := setup(t, common)

	t.Run("Nil ID", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		_, _, err := repo.ChangeChannelSubscription(uuid.Nil, ChangeChannelSubscriptionArgs{})
		assert.EqualError(err, ErrNilID.Error())
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		ch := mustMakeChannel(t, repo, rand)
		user1 := mustMakeUser(t, repo, rand)
		user2 := mustMakeUser(t, repo, rand)

		args := ChangeChannelSubscriptionArgs{
			Subscription: map[uuid.UUID]model.ChannelSubscribeLevel{
				user1.GetID():           model.ChannelSubscribeLevelMarkAndNotify,
				user2.GetID():           model.ChannelSubscribeLevelMarkAndNotify,
				uuid.Must(uuid.NewV4()): model.ChannelSubscribeLevelMarkAndNotify,
			},
		}
		_, _, err := repo.ChangeChannelSubscription(ch.ID, args)
		if assert.NoError(err) {
			assert.Equal(2, count(t, getDB(repo).Model(model.UserSubscribeChannel{}).Where(&model.UserSubscribeChannel{ChannelID: ch.ID})))
		}

		args = ChangeChannelSubscriptionArgs{
			Subscription: map[uuid.UUID]model.ChannelSubscribeLevel{
				user1.GetID():           model.ChannelSubscribeLevelMarkAndNotify,
				user2.GetID():           model.ChannelSubscribeLevelNone,
				uuid.Must(uuid.NewV4()): model.ChannelSubscribeLevelNone,
			},
		}
		_, _, err = repo.ChangeChannelSubscription(ch.ID, args)
		if assert.NoError(err) {
			assert.Equal(1, count(t, getDB(repo).Model(model.UserSubscribeChannel{}).Where(&model.UserSubscribeChannel{ChannelID: ch.ID})))
		}
	})
}

func TestGormRepository_GetChannelStats(t *testing.T) {
	t.Parallel()
	repo, _, _, user, channel := setupWithUserAndChannel(t, common)

	t.Run("nil id", func(t *testing.T) {
		t.Parallel()

		_, err := repo.GetChannelStats(uuid.Nil)
		assert.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		_, err := repo.GetChannelStats(uuid.Must(uuid.NewV4()))
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		for i := 0; i < 14; i++ {
			mustMakeMessage(t, repo, user.GetID(), channel.ID)
		}
		require.NoError(t, repo.DeleteMessage(mustMakeMessage(t, repo, user.GetID(), channel.ID).ID))

		stats, err := repo.GetChannelStats(channel.ID)
		if assert.NoError(t, err) {
			assert.NotEmpty(t, stats.DateTime)
			assert.EqualValues(t, 15, stats.TotalMessageCount)
		}
	})
}
