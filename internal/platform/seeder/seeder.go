package seeder

import (
	"blog/internal/platform/pkg/rbac"
	"github.com/casbin/casbin/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type Seeder struct {
	db       *mongo.Database
	enforcer *casbin.Enforcer
}

func NewSeeder(db *mongo.Database, e *casbin.Enforcer) *Seeder {
	return &Seeder{
		db:       db,
		enforcer: e,
	}
}

func (s *Seeder) createAdminPolicy() error {
	adminPolicy := []struct {
		role rbac.Role
		obj  rbac.Object
		act  rbac.Action
	}{
		{rbac.RoleAdmin, rbac.ObjectUser, rbac.ActionRead},
		{rbac.RoleAdmin, rbac.ObjectUser, rbac.ActionWrite},
		{rbac.RoleAdmin, rbac.ObjectUser, rbac.ActionModify},
		{rbac.RoleAdmin, rbac.ObjectUser, rbac.ActionDelete},
		{rbac.RoleAdmin, rbac.ObjectPost, rbac.ActionRead},
		{rbac.RoleAdmin, rbac.ObjectPost, rbac.ActionWrite},
		{rbac.RoleAdmin, rbac.ObjectPost, rbac.ActionModify},
		{rbac.RoleAdmin, rbac.ObjectPost, rbac.ActionDelete},
	}
	for _, policy := range adminPolicy {
		if err := s.ensurePolicy(policy.role.String(), policy.obj.String(), policy.act.String()); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seeder) createManagerPolicy() error {
	adminPolicy := []struct {
		role rbac.Role
		obj  rbac.Object
		act  rbac.Action
	}{
		{rbac.RoleManager, rbac.ObjectUser, rbac.ActionRead},
		{rbac.RoleManager, rbac.ObjectUser, rbac.ActionModify},
		{rbac.RoleManager, rbac.ObjectPost, rbac.ActionRead},
		{rbac.RoleManager, rbac.ObjectPost, rbac.ActionWrite},
		{rbac.RoleManager, rbac.ObjectPost, rbac.ActionModify},
	}
	for _, policy := range adminPolicy {
		if err := s.ensurePolicy(policy.role.String(), policy.obj.String(), policy.act.String()); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seeder) createEditorPolicy() error {
	adminPolicy := []struct {
		role rbac.Role
		obj  rbac.Object
		act  rbac.Action
	}{
		{rbac.RoleEditor, rbac.ObjectPost, rbac.ActionRead},
		{rbac.RoleEditor, rbac.ObjectPost, rbac.ActionWrite},
		{rbac.RoleEditor, rbac.ObjectPost, rbac.ActionModify},
	}
	for _, policy := range adminPolicy {
		if err := s.ensurePolicy(policy.role.String(), policy.obj.String(), policy.act.String()); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seeder) createAuthorPolicy() error {
	adminPolicy := []struct {
		role rbac.Role
		obj  rbac.Object
		act  rbac.Action
	}{
		{rbac.RoleAuthor, rbac.ObjectPost, rbac.ActionRead},
		{rbac.RoleAuthor, rbac.ObjectPost, rbac.ActionWrite},
		{rbac.RoleAuthor, rbac.ObjectPost, rbac.ActionModify},
		{rbac.RoleAuthor, rbac.ObjectPost, rbac.ActionDelete},
	}
	for _, policy := range adminPolicy {
		if err := s.ensurePolicy(policy.role.String(), policy.obj.String(), policy.act.String()); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seeder) createModeratorPolicy() error {
	adminPolicy := []struct {
		role rbac.Role
		obj  rbac.Object
		act  rbac.Action
	}{
		{rbac.RoleModerator, rbac.ObjectPost, rbac.ActionRead},
		{rbac.RoleModerator, rbac.ObjectPost, rbac.ActionModify},
	}
	for _, policy := range adminPolicy {
		if err := s.ensurePolicy(policy.role.String(), policy.obj.String(), policy.act.String()); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seeder) createMemberPolicy() error {
	adminPolicy := []struct {
		role rbac.Role
		obj  rbac.Object
		act  rbac.Action
	}{
		{rbac.RoleMember, rbac.ObjectPost, rbac.ActionRead},
	}
	for _, policy := range adminPolicy {
		if err := s.ensurePolicy(policy.role.String(), policy.obj.String(), policy.act.String()); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seeder) createGuestPolicy() error {
	adminPolicy := []struct {
		role rbac.Role
		obj  rbac.Object
		act  rbac.Action
	}{
		{rbac.RoleGuest, rbac.ObjectPost, rbac.ActionRead},
	}
	for _, policy := range adminPolicy {
		if err := s.ensurePolicy(policy.role.String(), policy.obj.String(), policy.act.String()); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seeder) ensurePolicy(role, obj, act string) error {
	if hasPolicy, err := s.enforcer.HasPolicy(role, obj, act); !hasPolicy {
		_, err = s.enforcer.AddPolicy(role, obj, act)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Seeder) Seed() error {
	err := s.createAdminPolicy()
	if err != nil {
		return err
	}
	err = s.createManagerPolicy()
	if err != nil {
		return err
	}
	err = s.createEditorPolicy()
	if err != nil {
		return err
	}
	err = s.createAuthorPolicy()
	if err != nil {
		return err
	}
	err = s.createModeratorPolicy()
	if err != nil {
		return err
	}
	err = s.createMemberPolicy()
	if err != nil {
		return err
	}
	err = s.createGuestPolicy()
	if err != nil {
		return err
	}
	return err
}
