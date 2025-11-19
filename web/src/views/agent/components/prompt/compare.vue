<template>
  <div class="compare-container">
    <div class="compare-header">
      <div class="compare-header-left">
        <span class="el-icon-arrow-left go-back" @click="goBack"></span>
        <h3>{{ $t('tempSquare.promptCompare') }}</h3>
      </div>
      <el-button type="primary" size="small" @click="addPromptField">
        <i class="el-icon-plus" style="margin-right:4px;"></i>
        {{ $t('tempSquare.addPrompt') }}
      </el-button>
    </div>
    <div class="compare-content">
      <div class="compare-middle">
        <div class="prompt-field-list">
          <div
            class="prompt-field-item"
            :style="{ width: `calc((100% - ${(promptFields.length - 1) * 10}px) / ${promptFields.length})` }"
            v-for="(field, index) in promptFields"
            :key="field.id"
          >
            <PromptField 
            :ref="getPromptFieldRef(index)"
            :fieldIndex="index" 
            :editForm="editForm" 
            @closePrompt="closePrompt"
            />
          </div>
        </div>
      </div>
      <div class="compare-bottom">
        <EditableDivV3
          ref="editable"
          source="promptCompare"
          :fileTypeArr="fileTypeArr"
          type="compare"
          @preSend="handlePromptSubmit"
        />
      </div>
    </div>
  </div>
</template>

<script>
import PromptField from './promptField.vue'
import EditableDivV3 from '../EditableDivV3'
import {getAgentDetail} from "@/api/agent";

export default {
  name: 'PromptCompare',
  components: {
    PromptField,
    EditableDivV3
  },
  data() {
    return {
      promptFields: [{ id: Date.now()+Math.random() }],
      fileTypeArr: [],
      currentModel: null,
      editForm:null
    }
  },
  created() {
    this.getAgentDetail()
  },
  methods: {
    handlePromptSubmit(){
      const editable = this.$refs.editable
      if (!editable || typeof editable.getValue !== 'function') return
      const promptText = (editable.getValue() || '').trim()
      if (!promptText) return

      for (let i = 0; i < this.promptFields.length; i += 1) {
        const refName = this.getPromptFieldRef(i)
        const fieldRef = this.$refs[refName]
        if (fieldRef && typeof fieldRef.runPrompt === 'function') {
          fieldRef.runPrompt(promptText)
        }
      }
    },
    getPromptFieldRef(index){
      return 'promptField' + index;
    },
    getAgentDetail() {
      getAgentDetail({ assistantId: this.$route.params.id }).then(res => {
        if (res.code === 0) {
          const detail = res.data || {}
          const recommendQuestion = Array.isArray(detail.recommendQuestion)
            ? detail.recommendQuestion.filter(item => item && item.value)
            : []
          this.editForm = {
            ...detail,
            recommendQuestion
          }
        }
      })     
    },
    closePrompt(index) {
      this.promptFields.splice(index, 1)
    },
    addPromptField() {
      if (this.promptFields.length > 4) {
        this.$message.warning(this.$t('tempSquare.promptCompareLimit'))
        return
      }
      this.promptFields.push({ id: Date.now() + Math.random() })
    },
    goBack(){
      const id = this.$route.params.id
      this.$router.push({ path: `/agent/test/?id=${id}` });
    }
  }
}
</script>

<style scoped lang="scss">
.compare-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
  padding:0 10px;
  box-sizing: border-box;
}

.compare-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding:24px 10px 14px 10px;
  border-bottom: 1px solid #eaeaea;
  .compare-header-left{
    display: flex;
    align-items: center;
    .go-back{
      font-size: 18px;
      cursor: pointer;
    }
    h3{
      font-size: 18px;
      font-weight: 800;
      color: #434c6c;
      margin-left: 10px;
  } 
  }
}
.compare-content{
  flex: 1;
  display:flex;
  flex-direction: column;
  box-sizing: border-box;
  .compare-middle {
    flex: 1;
    margin: 10px 10px 0 10px;
    overflow-y:auto;
    .prompt-field-list {
      display: flex;
      height: 100%;
      gap: 10px;
      overflow-x: auto;
      .prompt-field-item {
        flex: 1;
        height: 100%;
      }
    }
  }
  .compare-bottom {
    padding:10px;
  }
}

</style>
