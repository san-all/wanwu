package orm

import (
	"context"
	"errors"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/model"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/orm/sqlopt"
	"gorm.io/gorm"
)

func (c *Client) CreateAssistantSnapshot(ctx context.Context, assistantSnapshot *model.AssistantSnapshot) (uint32, *err_code.Status) {
	return assistantSnapshot.ID, c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		// 检查是否已存在相同的记录
		if err := sqlopt.SQLOptions(
			sqlopt.WithAssistantID(assistantSnapshot.AssistantID),
			sqlopt.WithVersion(assistantSnapshot.Version),
			sqlopt.DataPerm(assistantSnapshot.UserId, assistantSnapshot.OrgId),
		).Apply(tx).First(&model.AssistantSnapshot{}).Error; err == nil {
			return toErrStatus("assistant_snapshot", "snapshot already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return toErrStatus("assistant_snapshot", err.Error())
		}

		// 创建
		if err := tx.Create(assistantSnapshot).Error; err != nil {
			return toErrStatus("assistant_snapshot", err.Error())
		}
		return nil
	})
}

func (c *Client) UpdateAssistantSnapshot(ctx context.Context, assistantID uint32, desc string, userID, orgID string) *err_code.Status {
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		// 查询最新版本号
		var id uint32
		err := sqlopt.SQLOptions(
			sqlopt.WithAssistantID(assistantID),
			sqlopt.DataPerm(userID, orgID),
		).Apply(tx).Model(&model.AssistantSnapshot{}).
			Order("created_at DESC").
			Limit(1).
			Pluck("id", &id).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return toErrStatus("assistant_snapshot", "snapshot not found")
			}
			return toErrStatus("assistant_snapshot", "query snapshot version failed: "+err.Error())
		}

		if id == 0 {
			return toErrStatus("assistant_snapshot", "snapshot not found")
		}

		// 更新
		result := sqlopt.WithID(id).Apply(tx).Model(&model.AssistantSnapshot{}).Updates(map[string]interface{}{
			"desc": desc,
		})
		if result.Error != nil {
			return toErrStatus("assistant_snapshot", "update snapshot failed: "+result.Error.Error())
		}

		return nil
	})
}

func (c *Client) GetAssistantSnapshotList(ctx context.Context, assistantID uint32, userID, orgID string) ([]*model.AssistantSnapshot, *err_code.Status) {
	var assistantSnapshots []*model.AssistantSnapshot
	err := sqlopt.SQLOptions(
		sqlopt.WithAssistantID(assistantID),
		sqlopt.DataPerm(userID, orgID),
	).Apply(c.db.WithContext(ctx)).Model(&model.AssistantSnapshot{}).
		Order("created_at DESC").
		Find(&assistantSnapshots).Error
	if err != nil {
		return nil, toErrStatus("assistant_snapshot_list", err.Error())
	}
	return assistantSnapshots, nil
}

func (c *Client) GetAssistantSnapshot(ctx context.Context, assistantID uint32, version string, userID, orgID string) (*model.AssistantSnapshot, *err_code.Status) {
	assistantSnapshot := &model.AssistantSnapshot{}
	err := sqlopt.SQLOptions(
		sqlopt.WithAssistantID(assistantID),
		sqlopt.WithVersion(version),
		sqlopt.DataPerm(userID, orgID),
	).Apply(c.db.WithContext(ctx)).Model(&model.AssistantSnapshot{}).
		Order("created_at DESC").
		First(&assistantSnapshot).Error
	if err != nil {
		return nil, toErrStatus("assistant_snapshot", err.Error())
	}
	return assistantSnapshot, nil
}

func (c *Client) RollbackAssistantSnapshot(ctx context.Context, assistant *model.Assistant, tools []*model.AssistantTool, mcps []*model.AssistantMCP, workflows []*model.AssistantWorkflow, userID, orgID string) *err_code.Status {
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		// Update Assistant Info
		if assistant.ID == 0 {
			return toErrStatus("assistant_update", "update assistant but id 0")
		}
		if err := tx.Model(assistant).Updates(map[string]interface{}{
			"avatar_path":          assistant.AvatarPath,
			"name":                 assistant.Name,
			"desc":                 assistant.Desc,
			"instructions":         assistant.Instructions,
			"prologue":             assistant.Prologue,
			"recommend_question":   assistant.RecommendQuestion,
			"model_config":         assistant.ModelConfig,
			"knowledgebase_config": assistant.KnowledgebaseConfig,
			"scope":                assistant.Scope,
			"rerank_config":        assistant.RerankConfig,
			"safety_config":        assistant.SafetyConfig,
			"vision_config":        assistant.VisionConfig,
		}).Error; err != nil {
			return toErrStatus("assistant_update", err.Error())
		}

		assistantId := assistant.ID

		// Snapshot Tools
		// 删除该智能体所绑定的所有工具
		if err := sqlopt.SQLOptions(
			sqlopt.WithAssistantID(assistantId),
			sqlopt.DataPerm(userID, orgID),
		).Apply(tx).Delete(&model.AssistantTool{}).Error; err != nil {
			return toErrStatus("assistant_tool", err.Error())
		}
		// 创建
		if len(tools) > 0 {
			if err := tx.Create(tools).Error; err != nil {
				return toErrStatus("assistant_tool", err.Error())
			}
		}

		// Snapshot MCPs
		// 删除该智能体所绑定的所有mcp
		if err := sqlopt.SQLOptions(
			sqlopt.WithAssistantID(assistantId),
			sqlopt.DataPerm(userID, orgID),
		).Apply(tx).Delete(&model.AssistantMCP{}).Error; err != nil {
			return toErrStatus("assistant_mcp", err.Error())
		}
		// 创建
		if len(mcps) > 0 {
			if err := tx.Create(mcps).Error; err != nil {
				return toErrStatus("assistant_mcp", err.Error())
			}
		}

		// Snapshot Workflows
		// 删除该智能体所绑定的所有workflow
		if err := sqlopt.SQLOptions(
			sqlopt.WithAssistantID(assistantId),
			sqlopt.DataPerm(userID, orgID),
		).Apply(tx).Delete(&model.AssistantWorkflow{}).Error; err != nil {
			return toErrStatus("assistant_workflow", err.Error())
		}
		// 创建
		if len(workflows) > 0 {
			if err := tx.Create(workflows).Error; err != nil {
				return toErrStatus("assistant_workflow", err.Error())
			}
		}

		return nil
	})
}
