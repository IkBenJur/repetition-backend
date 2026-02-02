package types

import "time"

// Represent common payload base for all datatables
type BasePayloadEntity struct {
	Id        *int64     `json:"id"`
	UserId    *int64     `json:"user_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (entity *BasePayloadEntity) ToBaseEntity() *BaseEntity {
	base := &BaseEntity{}

	// Copy Id when set
	if entity.Id != nil {
		base.Id = entity.Id
	}

	// Copy userId when set
	if entity.UserId != nil {
		base.UserId = entity.UserId
	}

	if entity.CreatedAt != nil {
		base.CreatedAt = entity.CreatedAt
	}

	if entity.UpdatedAt != nil {
		base.UpdatedAt = entity.UpdatedAt
	}

	return base
}

func (entity *BasePayloadEntity) IsUpdate() bool {
	return entity.Id != nil && *entity.Id > 0
}

func (entity *BasePayloadEntity) IsNew() bool {
	return entity.Id == nil
}

// Represent common base for all datatables
type BaseEntity struct {
	Id        *int64
	UserId    *int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (b *BaseEntity) IsNew() bool {
	return b.Id == nil
}

func (b *BaseEntity) IsUpdate() bool {
	return b.Id != nil && *b.Id > 0
}

func (b *BaseEntity) MarkCreated() {
	now := time.Now()
	b.CreatedAt = &now
	b.UpdatedAt = &now
}

func (b *BaseEntity) MarkUpdated() {
	now := time.Now()
	b.UpdatedAt = &now
}

// This function can be used to correctly update the right timestamp
func (b *BaseEntity) Touch() {
	if b.IsNew() {
		b.MarkCreated()
	} else {
		b.MarkUpdated()
	}
}
