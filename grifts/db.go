package grifts

import (
	"fmt"
	"strconv"

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

			// users and roles seeding
			usersAndRoles := []interface{}{
				userAdmin,
				userMember,
				userAnonymous,
			}
			for _, m := range usersAndRoles {
				if err := createModel(m, tx); err != nil {
					fmt.Printf("%v", err)
					return err
				}
			}

			// products seeding
			for i := 0; i < 20; i++ {
				p := models.Product{
					Name:     "Product " + strconv.Itoa(i),
					ImageURL: "https://picsum.photos/350/250?image=" + strconv.Itoa(i),
					Price:    17,
				}
				if err := createModel(&p, tx); err != nil {
					fmt.Printf("%v", err)
					return err
				}
			}

			return nil
		})

		return nil
	})

})

func createModel(r interface{}, tx *pop.Connection) error {
	verrs, err := tx.Eager().ValidateAndCreate(r)
	if err != nil || verrs.HasAny() {
		return errors.WithStack(err)
	}

	return nil
}
