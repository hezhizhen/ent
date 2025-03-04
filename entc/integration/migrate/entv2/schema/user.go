// Copyright 2019-present Facebook Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type Mixin struct {
	mixin.Schema
}

func (m Mixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("mixed_string").
			Default("default"),
		field.Enum("mixed_enum").
			Values("on", "off").
			Default("on"),
	}
}

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Mixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			StorageKey("oid"),
		// changing the type of the field.
		field.Int("age"),
		// extending name field to longtext.
		field.Text("name"),
		// extending the index prefix below (on MySQL).
		field.Text("description").
			Optional(),
		// changing nickname from unique no non-unique.
		field.String("nickname").
			MaxLen(255),
		// adding new columns (must be either optional, or with a default value).
		field.String("phone").
			Default("unknown"),
		field.Bytes("buffer").
			Optional().
			DefaultFunc(func() []byte { return []byte("null") }),
		// adding new column with supported default value
		// in the database side, will append this value to
		// all existing rows.
		field.String("title").
			Default("SWE"),
		// change column name and reference it to the
		// previous one ("renamed").
		field.String("new_name").
			Optional().
			StorageKey("renamed"),
		// extending the blob size.
		field.Bytes("blob").
			Optional().
			MaxLen(1000),
		// adding enum to the `state` column.
		field.Enum("state").
			Optional().
			Values("logged_in", "logged_out", "online"),
		// convert string to enum.
		field.Enum("status").
			Optional().
			Values("done", "pending"),
		// remove the max-length constraint from varchar.
		field.String("workplace").
			Optional(),
		// add a new column with generated values by the database.
		field.Time("created_at").
			Default(time.Now).
			Annotations(&entsql.Annotation{
				Default: "CURRENT_TIMESTAMP",
			}),
		// deleting the `address` column.
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge(children<-M2O->parent) to be dropped.
		// Edge(spouse<-O2O->spouse) to be dropped.
		edge.To("car", Car.Type),
		// New edges to added.
		edge.To("pets", Pet.Type).
			StorageKey(edge.Column("owner_id"), edge.Symbol("user_pet_id")).
			Unique(),
		edge.To("friends", User.Type).
			StorageKey(
				edge.Table("friends"),
				edge.Columns("user", "friend"),
				edge.Symbols("user_friend_id1", "user_friend_id2"),
			),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		// Extend the column prefix by drop and create
		// this index on MySQL.
		index.Fields("description").
			Annotations(entsql.Prefix(100)),
		// Deleting old indexes (name, address),
		// and defining a new one.
		index.Fields("phone", "age").
			Unique(),
	}
}

type Car struct {
	ent.Schema
}

func (Car) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("car").
			Unique().
			// Make a M20 edge from nullable to required.
			// Requires column and foreign-key migration.
			Required(),
	}
}

// Group schema.
type Group struct{ ent.Schema }

// Pet schema.
type Pet struct {
	ent.Schema
}

func (Pet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("pets").
			Unique(),
	}
}
