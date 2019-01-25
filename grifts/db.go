package grifts

import (
	"github.com/gobuffalo/pop"
	"github.com/kgosse/shop-back/models"
	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		DB, err := pop.Connect("development")
		if err != nil {
			return err
		}

		DB.Transaction(func(tx *pop.Connection) error {
			// default roles
			roleAdmin := models.Role{Role: "admin"}
			roleMember := models.Role{Role: "member"}
			roleAnonymous := models.Role{Role: "anonymous"}
			// default users
			userAdmin := &models.User{
				Email:    "admin@gmail.com",
				Name:     "admin",
				Password: "admin",
				Roles:    models.Roles{roleAdmin},
			}
			userMember := &models.User{
				Email:    "member@gmail.com",
				Name:     "member",
				Password: "member",
				Roles:    models.Roles{roleMember},
			}
			userAnonymous := &models.User{
				Email:    "anonymous@gmail.com",
				Name:     "anonymous",
				Password: "anonymous",
				Roles:    models.Roles{roleAnonymous},
			}
			// @todo
			// default products

			// db seeding
			models := []interface{}{
				userAdmin,
				userMember,
				userAnonymous,
			}
			for _, m := range models {
				if err := createModel(m, tx); err != nil {
					return err
				}
			}

			return nil
		})

		return nil
	})

})

func createModel(r interface{}, tx *pop.Connection) error {
	err := tx.Eager().Create(r)
	if err != nil {
		return errors.WithStack(err)
	}

	// if verrs.HasAny() {
	// 	return errors.WithStack(verrs)
	// }
	return nil
}
