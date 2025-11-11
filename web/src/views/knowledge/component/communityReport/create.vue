<template>
  <el-dialog
    :title="$t('knowledgeManage.communityReport.addCommunityReport')"
    :visible.sync="dialogVisible"
    width="50%"
    :before-close="handleClose"
  >
    <el-form
      :model="ruleForm"
      ref="ruleForm"
      label-width="120px"
      class="demo-ruleForm"
    >
      <el-form-item
        class="itemCenter"
      >
        <el-radio-group
          v-model="createType"
          @input="typeChange($event)"
        >
          <el-radio-button :label="'single'">{{$t('knowledgeManage.create.single')}}</el-radio-button>
          <el-radio-button :label="'file'">{{$t('knowledgeManage.create.file')}}</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item
        :label="$t('knowledgeManage.create.file')"
        v-if="createType === 'file'"
        prop="fileUploadId"
        :rules="[{ required: true, message: $t('common.input.placeholder'), trigger: 'blur' }]"
      >
        <fileUpload
          ref="fileUpload"
          :templateUrl="templateUrl"
          @uploadFile="uploadFile"
          :accept="accept"
        />
      </el-form-item>
      <template v-if="createType === 'single'">
        <el-form-item
          :label="$t('knowledgeManage.create.title')"
          prop="content"
          :rules="[{ required: true, message: $t('knowledgeManage.create.titlePlaceholder'), trigger: 'blur' }]"
        >
          <el-input
            :placeholder="$t('knowledgeManage.create.titlePlaceholder')"
            v-model="ruleForm.title"
            maxlength="100"
            show-word-limit
          ></el-input>
        </el-form-item>
        <el-form-item
          :label="$t('knowledgeManage.create.content')"
          prop="content"
          :rules="[{ required: true, message: $t('knowledgeManage.create.contentPlaceholder'), trigger: 'blur' }]"
        >
          <el-input
            :placeholder="$t('knowledgeManage.create.contentPlaceholder')"
            v-model="ruleForm.content"
            type="textarea"
            :rows="6"
          ></el-input>
        </el-form-item>
        <el-form-item :label="$t('knowledgeManage.create.typeTitle')">
          <el-checkbox-group v-model="checkType">
            <el-checkbox
              label="more"
              name="type"
            >{{$t('knowledgeManage.create.continue')}}</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </template>
    </el-form>
    <span
      slot="footer"
      class="dialog-footer"
    >
      <el-button @click="dialogVisible = false">{{ $t('common.confirm.cancel') }}</el-button>
      <el-button
        type="primary"
        @click="submit('ruleForm')"
        :loading="btnLoading"
      >{{ $t('common.confirm.confirm') }}</el-button>
    </span>
  </el-dialog>
</template>
<script>
import fileUpload from "@/components/fileUpload";
import {
  createSegment,
  createBatchSegment,
  createSegmentChild,
} from "@/api/knowledge";
export default {
  components: { fileUpload },
  props: {
    parentId: {
      type: String,
      default: ""
    }
  },
  data() {
    return {
      btnLoading: false,
      accept: ".csv",
      checkType: [],
      inputVisible: false,
      inputValue: "",
      createType: "single",
      ruleForm: {
        title: "",
        docId: "",
        fileUploadId: "",
      },
      dialogVisible: false,
      templateUrl: "/user/api/v1/static/docs/segment.csv",
      isChildChunk: false,
    };
  },
  methods: {
    typeChange(val) {
      if (val === "single") {
        this.ruleForm.fileUploadId = "";
        this.$refs.fileUpload && this.$refs.fileUpload.clearFileList();
      } else {
        this.clearForm();
        this.$refs.ruleForm && this.$refs.ruleForm.clearValidate();
      }
    },
    uploadFile(fileUploadId) {
      this.ruleForm.fileUploadId = fileUploadId;
    },
    handleClose() {
      this.dialogVisible = false;
    },
    showDialog(docId) {
      this.dialogVisible = true;
      this.ruleForm.docId = docId;
      this.clearForm();
    },
    submit(formName) {
      if (this.createType === "single") {
        this.handleSingle(formName);
      } else {
        this.handleFile();
      }
    },
    handleSingle(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.btnLoading = true;
          if (this.isChildChunk) {
            this.createChildChunk();
          } else {
            this.createParentChunk();
          }
        } else {
          return false;
        }
      });
    },
    createParentChunk() {
      const data = this.isChildChunk
        ? { content: this.ruleForm.content, docId: this.ruleForm.docId }
        : {
            content: this.ruleForm.content,
            docId: this.ruleForm.docId,
            labels: this.ruleForm.labels,
          };
      createSegment(data)
        .then((res) => {
          if (res.code === 0) {
            this.$message.success(this.$t('knowledgeManage.create.createSuccess'));
            if (!this.checkType.length) {
              this.dialogVisible = false;
              this.$emit("updateDataBatch");
            } else {
              this.clearForm();
              this.$emit("updateData");
            }
            this.btnLoading = false;
          }
        })
        .catch(() => {
          this.btnLoading = false;
        });
    },
    createChildChunk() {
      const data = {
        content: [this.ruleForm.content],
        docId: this.ruleForm.docId,
        parentId: this.parentId
      };
      createSegmentChild(data).then((res) => {
        if (res.code === 0) {
            this.$message.success(this.$t('knowledgeManage.create.createSuccess'));
            if (!this.checkType.length) {
              this.dialogVisible = false;
            } else {
              this.clearForm();
            }
            this.$emit("updateChildData");
            this.btnLoading = false;
          }
      }).catch(() => {
        this.btnLoading = false;
      });
    },
    handleFile() {
      this.btnLoading = true;
      const data = {
        fileUploadId: this.ruleForm.fileUploadId,
        docId: this.ruleForm.docId,
      };
      createBatchSegment(data)
        .then((res) => {
          if (res.code === 0) {
            this.$message.success(this.$t('knowledgeManage.create.createSuccess'));
            this.dialogVisible = false;
            this.btnLoading = false;
            this.$emit("updateDataBatch");
          }
        })
        .catch(() => {
          this.btnLoading = false;
        });
    },
    clearForm() {
      this.ruleForm.content = "";
      this.ruleForm.title = "";
      this.ruleForm.fileUploadId = "";
      this.checkType = [];
    },
  },
};
</script>
<style lang="scss" scoped>
.itemCenter {
  display: flex;
  justify-content: center;
  /deep/.el-form-item__content {
    margin-left: 0 !important;
  }
}
.el-tag {
  margin-right: 5px;
  color: #3848f7;
  border-color: #3848f7;
  background: $color_opacity;
}
/deep/ {
  .el-tag .el-tag__close {
    color: #3848f7 !important;
  }
  .el-tag .el-tag__close:hover {
    color: #fff !important;
    background: #3848f7;
  }
  .el-checkbox__input.is-checked + .el-checkbox__label {
    color: #3848f7;
  }
}
</style>
