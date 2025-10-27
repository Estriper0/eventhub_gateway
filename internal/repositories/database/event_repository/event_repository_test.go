package event_repository

import (
	"context"
	"testing"
	"time"

	"github.com/Estriper0/EventHub/internal/config"
	"github.com/Estriper0/EventHub/internal/models"
	"github.com/Estriper0/EventHub/internal/repositories"
	"github.com/Estriper0/EventHub/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEventRepo_Create(t *testing.T) {
	t.Setenv("ENV", "test")
	config := config.New()
	db, teardown := testutils.GetDb(t, &config.DB)
	defer teardown()

	repo := New(db)
	ctx := context.Background()

	event_test := &models.EventCreateRequest{
		Title:        "Test_Title",
		About:        "Test_About",
		StartDate:    time.Now(),
		Location:     "Test_Location",
		Status:       "Test_Status",
		MaxAttendees: 15,
	}

	id, err := repo.Create(ctx, event_test)

	assert.NoError(t, err)
	assert.NotNil(t, id)
}

func TestEventRepo_GetById(t *testing.T) {
	t.Setenv("ENV", "test")
	config := config.New()
	db, teardown := testutils.GetDb(t, &config.DB)
	defer teardown()

	repo := New(db)
	ctx := context.Background()

	event_test := &models.EventCreateRequest{
		Title:        "Test_Title",
		About:        "Test_About",
		StartDate:    time.Now().UTC(),
		Location:     "Test_Location",
		Status:       models.StatusCompleted,
		MaxAttendees: 15,
	}

	id, err := repo.Create(ctx, event_test)

	assert.NoError(t, err)
	assert.NotNil(t, id)

	event, err := repo.GetById(ctx, id)

	assert.NoError(t, err)
	assert.Equal(t, event.Title, event_test.Title)
	assert.Equal(t, event.About, event_test.About)
	assert.Equal(t, event.StartDate.Sub(event_test.StartDate).Abs() < 1*time.Millisecond, true)
	assert.Equal(t, event.Location, event_test.Location)
	assert.Equal(t, event.Status, event_test.Status)
	assert.Equal(t, event.MaxAttendees, event_test.MaxAttendees)

	event, err = repo.GetById(ctx, uuid.New())
	assert.EqualError(t, err, repositories.ErrRecordNotFound.Error())
	assert.Nil(t, event)
}

func TestEventRepo_GetAll(t *testing.T) {
	t.Setenv("ENV", "test")
	config := config.New()
	db, teardown := testutils.GetDb(t, &config.DB)
	defer teardown()

	repo := New(db)
	ctx := context.Background()

	events, err := repo.GetAll(ctx)

	assert.NoError(t, err)
	assert.Equal(t, len(events), 0)

	events_test := []*models.EventCreateRequest{
		{
			Title:        "Test_Title1",
			About:        "Test_About1",
			StartDate:    time.Now().UTC(),
			Location:     "Test_Location1",
			Status:       models.StatusCancelled,
			MaxAttendees: 15,
		},
		{
			Title:        "Test_Title2",
			About:        "Test_About2",
			StartDate:    time.Now().UTC(),
			Location:     "Test_Location2",
			Status:       models.StatusCompleted,
			MaxAttendees: 20,
		},
		{
			Title:        "Test_Title3",
			About:        "Test_About3",
			StartDate:    time.Now().UTC(),
			Location:     "Test_Location3",
			Status:       models.StatusDraft,
			MaxAttendees: 25,
		},
	}

	for _, ev := range events_test {
		_, err := repo.Create(ctx, ev)
		assert.NoError(t, err)
	}

	events, err = repo.GetAll(ctx)

	assert.NoError(t, err)
	assert.Equal(t, len(events), len(events_test))

	for i, ev := range events {
		assert.Equal(t, ev.Title, events_test[i].Title)
		assert.Equal(t, ev.About, events_test[i].About)
		assert.Equal(t, ev.StartDate.Sub(events_test[i].StartDate).Abs() < 1*time.Millisecond, true)
		assert.Equal(t, ev.Location, events_test[i].Location)
		assert.Equal(t, ev.Status, events_test[i].Status)
		assert.Equal(t, ev.MaxAttendees, events_test[i].MaxAttendees)
	}
}

func TestEventRepo_DeleteById(t *testing.T) {
	t.Setenv("ENV", "test")
	config := config.New()
	db, teardown := testutils.GetDb(t, &config.DB)
	defer teardown()

	repo := New(db)
	ctx := context.Background()

	event_test := &models.EventCreateRequest{
		Title:        "Test_Title",
		About:        "Test_About",
		StartDate:    time.Now().UTC(),
		Location:     "Test_Location",
		Status:       models.StatusDraft,
		MaxAttendees: 15,
	}

	id, err := repo.Create(ctx, event_test)

	assert.NoError(t, err)
	assert.NotNil(t, id)

	err = repo.DeleteById(ctx, id)

	assert.NoError(t, err)

	events, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, len(events), 0)

	err = repo.DeleteById(ctx, id)
	assert.EqualError(t, err, repositories.ErrRecordNotFound.Error())
}

func TestEventRepo_Update(t *testing.T) {
	t.Setenv("ENV", "test")
	config := config.New()
	db, teardown := testutils.GetDb(t, &config.DB)
	defer teardown()

	repo := New(db)
	ctx := context.Background()

	event_test := &models.EventCreateRequest{
		Title:        "Test_Title",
		About:        "Test_About",
		StartDate:    time.Now().UTC(),
		Location:     "Test_Location",
		Status:       models.StatusCancelled,
		MaxAttendees: 15,
	}

	id, err := repo.Create(ctx, event_test)

	assert.NoError(t, err)
	assert.NotNil(t, id)

	new_title := "Test_Update"
	new_max_attendess := 15
	event_update := &models.EventUpdateRequest{
		Id:           &id,
		Title:        &new_title,
		MaxAttendees: &new_max_attendess,
	}

	err = repo.Update(ctx, event_update)
	assert.NoError(t, err)

	event, err := repo.GetById(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, event.Title, *event_update.Title)
	assert.Equal(t, event.MaxAttendees, *event_update.MaxAttendees)

	*event_update.Id = uuid.New()
	err = repo.Update(ctx, event_update)
	assert.EqualError(t, err, repositories.ErrRecordNotFound.Error())
}
