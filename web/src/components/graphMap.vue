<template>
  <div class="graph-map-container">
    <div class="graph-map-content">
      <div class="graph-header" v-if="showHeader">
        <span class="el-icon-arrow-left back" @click="goBack"></span>
        <span class="header-text">
          {{headerName}}
        </span>
      </div>
      <div class="graph-toolbar">
        <el-tooltip :content="$t('knowledgeManage.graph.zoomIn')" placement="top">
          <el-button 
            icon="el-icon-zoom-in" 
            circle 
            size="mini" 
            @click="zoomIn"
          ></el-button>
        </el-tooltip>
        <el-tooltip :content="$t('knowledgeManage.graph.zoomOut')" placement="top">
          <el-button 
            icon="el-icon-zoom-out" 
            circle 
            size="mini" 
            @click="zoomOut"
          ></el-button>
        </el-tooltip>
        <el-tooltip :content="$t('knowledgeManage.graph.fitView')" placement="top">
          <el-button 
            icon="el-icon-full-screen" 
            circle 
            size="mini" 
            @click="fitView"
          ></el-button>
        </el-tooltip>
        <el-tooltip :content="$t('knowledgeManage.graph.resetZoom')" placement="top">
          <el-button 
            icon="el-icon-refresh-left" 
            circle 
            size="mini" 
            @click="resetZoom"
          ></el-button>
        </el-tooltip>
        <el-divider direction="vertical"></el-divider>
        <el-tooltip :content="$t('knowledgeManage.graph.refresh')" placement="top">
          <el-button 
            icon="el-icon-refresh" 
            circle 
            size="mini" 
            :loading="loading"
            @click="refreshData"
          ></el-button>
        </el-tooltip>
      </div>
      <div ref="graphContainer" class="graph-container" v-loading="loading"></div>
    </div>
  </div>
</template>

<script>
import G6 from '@antv/g6'

export default {
  name: 'GraphMap',
  props: {
    data: {
      type: Object,
      default: function() {
        return {
          nodes: [
            { id: '1', label: '中心节点', type: 'circle', size: 20, style: { fill: '#5B8FF9', stroke: '#1890ff' } },
            { id: '2', label: '节点A', type: 'circle', size: 20, style: { fill: '#C6E5FF' } },
            { id: '3', label: '节点B', type: 'circle', size: 20, style: { fill: '#C6E5FF' } },
            { id: '4', label: '节点C', type: 'circle', size: 20, style: { fill: '#C6E5FF' } },
            { id: '5', label: '节点D', type: 'circle', size: 20, style: { fill: '#C6E5FF' } }
          ],
          edges: [
            { id: 'e1', source: '1', target: '2', label: '连接1' },
            { id: 'e2', source: '1', target: '3', label: '连接2' },
            { id: 'e3', source: '1', target: '4', label: '连接3' },
            { id: 'e4', source: '1', target: '5', label: '连接4' },
            { id: 'e5', source: '2', target: '2', label: '连接5' },
          ]
        }
      }
    },
    showHeader:{
      type:Boolean,
      default:true
    },
    knowledgeId: {
      type: [String, Number],
      default: null
    },
    autoFit: {
      type: Boolean,
      default: true
    },
    layout: {
      type: String,
      default: 'force'
    },
    nodeStyle: {
      type: Object,
      default: function() {
        return {}
      }
    },
    edgeStyle: {
      type: Object,
      default: function() {
        return {}
      }
    },
    performance: {
      type: Object,
      default: function() {
        return {
          largeThreshold: 500,
          hideLabelZoom: 0.8,
          enableWorkerLayout: true,
          simplifyEdgeOnLarge: true,
          disableAnimateOnLarge: true
        }
      }
    }
  },
  data() {
    return {
      graph: null,
      loading: false,
      zoom: 1,
      lastLabelSwitchRAF: 0
    }
  },
  computed: {
    headerName() {
      if (this.$route && this.$route.query && this.$route.query.name) {
        return this.$route.query.name
      }
      return '知识库'
    }
  },
  watch: {
    data: {
      handler(newData) {
        if (this.graph && newData) {
          this.updateGraphData(newData)
        }
      },
      deep: true,
      immediate: false
    }
  },
  mounted() {
    this.initGraph()
  },
  beforeDestroy() {
    if (this.graph) {
      this.graph.destroy()
      this.graph = null
    }
    window.removeEventListener('resize', this.handleResize)
  },
  methods: {
    goBack(){
      let id = this.knowledgeId
      if (!id && this.$route && this.$route.params) {
        id = this.$route.params.id
      }
      
      if (id) {
        this.$router.push({
          path: `/knowledge/doclist/${id}`
        })
      }
    },
    initGraph() {
      if (!this.$refs.graphContainer) {
        return
      }

      const container = this.$refs.graphContainer
      const width = container.clientWidth || 800
      const height = container.clientHeight || 600

      this.registerCustomNodes()
      const Graph = G6.Graph || G6

      const graphConfig = {
        container,
        width,
        height,
        animate: false,
        modes: {
          default: [
            'drag-canvas',
            'zoom-canvas',
            'drag-node',
            'click-select',
            'brush-select'
          ]
        },
        layout: {
          type: this.layout,
          preventOverlap: true,
          nodeSize: 20,
          nodeSpacing: 50,
          ...this.getLayoutConfig()
        },
        defaultNode: {
          type: 'circle',
          size: 20,
          labelCfg: {
            position: 'bottom',
            offset: 10,
            style: {
              fill: '#333',
              fontSize: 14,
              fontWeight: 'normal',
              background: {
                fill: '#fff',
                padding: [2, 4, 2, 4],
                radius: 2
              }
            }
          },
          style: {
            fill: '#C6E5FF',
            stroke: '#5B8FF9',
            lineWidth: 2
          },
          ...this.nodeStyle
        },
        defaultEdge: {
          type: 'line',
          style: {
            stroke: '#A3B1BF',
            lineWidth: 2,
            endArrow: false
          },
          labelCfg: {
            autoRotate: true,
            refY: -10,
            style: {
              fill: '#666',
              fontSize: 14,
              fontWeight: 'normal',
              background: {
                fill: '#fff',
                padding: [2, 4, 2, 4],
                radius: 2,
                stroke: '#e4e7ed',
                lineWidth: 1
              }
            }
          },
          ...this.edgeStyle
        },
        nodeStateStyles: {
          hover: {
            fill: '#91d5ff',
            stroke: '#1890ff',
            lineWidth: 3
          },
          selected: {
            fill: '#91d5ff',
            stroke: '#1890ff',
            lineWidth: 3
          }
        },
        edgeStateStyles: {
          hover: {
            stroke: '#1890ff',
            lineWidth: 3
          },
          selected: {
            stroke: '#1890ff',
            lineWidth: 3
          }
        }
      }

      this.graph = new Graph(graphConfig)
 
      this.bindEvents()

      if (this.data && (this.data.nodes || this.data.edges)) {
        this.graph.data(this.data)
        this.graph.render()
        
        if (this.autoFit) {
          this.fitView()
        }
      }
 
      window.addEventListener('resize', this.handleResize)
    },

    registerCustomNodes() {
    },
 
    getLayoutConfig() {
      const configs = {
        force: {
          preventOverlap: true,
          nodeSize: 20,
          nodeSpacing: 50,
          linkDistance: 150,
          nodeStrength: -100,
          edgeStrength: 0.3,
          collideStrength: 0.8,
          alpha: 0.3,
          alphaDecay: 0.02,
          alphaMin: 0.001
        },
        dagre: {
          rankdir: 'TB',
          nodesep: 50,
          ranksep: 50
        },
        circular: {
          radius: 200,
          startRadius: 10,
          endRadius: 300
        },
        radial: {
          unitRadius: 100,
          nodeSize: 50
        }
      }
      return configs[this.layout] || configs.force
    },

    bindEvents() {
       if (!this.graph) return
 
       this.graph.on('node:click', (e) => {
         const node = e.item
         this.$emit('node-click', node.getModel())
       })
 
       this.graph.on('node:mouseenter', (e) => {
         const node = e.item
         this.graph.setItemState(node, 'hover', true)
       })
 
       this.graph.on('node:mouseleave', (e) => {
         const node = e.item
         this.graph.setItemState(node, 'hover', false)
       })
 
       this.graph.on('edge:click', (e) => {
         const edge = e.item
         this.$emit('edge-click', edge.getModel())
       })
 
       this.graph.on('canvas:click', () => {
         this.graph.getNodes().forEach(node => {
           this.graph.setItemState(node, 'selected', false)
         })
         this.graph.getEdges().forEach(edge => {
           this.graph.setItemState(edge, 'selected', false)
         })
       })
 
       this.graph.on('viewportchange', () => {
         if (this.graph) {
           const zoom = this.graph.getZoom()
          this.zoom = zoom
          this.$emit('zoom-change', zoom)
          this.toggleLabelsByZoom()
        }
      })
    },

    updateGraphData(data) {
      if (!this.graph) return
 
      const safeData = {
        nodes: Array.isArray(data && data.nodes)
          ? data.nodes.map(n => ({ ...n }))
          : [],
        edges: Array.isArray(data && data.edges)
          ? data.edges.map(e => {
              const { label, ...rest } = e || {}
              return { ...rest }
            })
          : []
      }
      this.graph.data(safeData)
      this.graph.render()
      
      if (this.autoFit) {
        this.$nextTick(() => {
          this.fitView()
        })
      }
    },
 
    toggleLabelsByZoom() {
      if (!this.graph) return
      const threshold = (this.performance && this.performance.hideLabelZoom) || 0
      if (threshold <= 0) return
      const now = (typeof performance !== 'undefined' && performance.now) ? performance.now() : Date.now()
      if (now - this.lastLabelSwitchRAF < 16) return
      this.lastLabelSwitchRAF = now
      const shouldHide = this.zoom < threshold
      this.graph.setAutoPaint(false)
      this.graph.getNodes().forEach((n) => {
        const model = n.getModel()
        const hasLabel = !!model.originalLabel
        const currentLabel = model.label
        if (shouldHide && currentLabel) {
          this.graph.updateItem(n, { label: '' })
        } else if (!shouldHide && !currentLabel && hasLabel) {
          this.graph.updateItem(n, { label: model.originalLabel })
        }
      })
      this.graph.setAutoPaint(true)
      this.graph.paint()
    },

    isLargeData(data) {
      const nodesLen = (data && data.nodes && data.nodes.length) || 0
      const threshold = (this.performance && this.performance.largeThreshold) || 500
      return nodesLen >= threshold
    },

    zoomIn() {
      if (!this.graph) return
      const currentZoom = this.graph.getZoom()
      const newZoom = Math.min(currentZoom * 1.2, 3)
      this.graph.zoomTo(newZoom)
    },
 
    zoomOut() {
      if (!this.graph) return
      const currentZoom = this.graph.getZoom()
      const newZoom = Math.max(currentZoom * 0.8, 0.3)
      this.graph.zoomTo(newZoom)
    },
 
    fitView() {
      if (!this.graph) return
      this.graph.fitView(20)
    },
 
    resetZoom() {
      if (!this.graph) return
      this.graph.zoomTo(1)
      this.fitView()
    },
 
    refreshData() {
      this.loading = true
      this.$emit('refresh', () => {
        this.loading = false
        if (this.autoFit) {
          this.$nextTick(() => {
            this.fitView()
          })
        }
      })
    },
     
    finishRefresh() {
      this.loading = false
      if (this.autoFit && this.graph) {
        this.$nextTick(() => {
          this.fitView()
        })
      }
    },
 
    handleResize() {
      if (!this.graph || !this.$refs.graphContainer) return
      
      const container = this.$refs.graphContainer
      const width = container.clientWidth
      const height = container.clientHeight
      
      this.graph.changeSize(width, height)
    },
 
    getGraph() {
      return this.graph
    },
 
    downloadImage(fileName = 'graph') {
      if (!this.graph) return
      
      if (typeof this.graph.downloadFullImage === 'function') {
        this.graph.downloadFullImage(fileName, 'image/png', {
          backgroundColor: '#fff',
          padding: [10, 10, 10, 10]
        })
      } else if (typeof this.graph.downloadImage === 'function') {
        this.graph.downloadImage(fileName)
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.graph-map-container {
  position: relative;
  width: 100%;
  height:100%;
  box-sizing: border-box;
  overflow: hidden;
  padding:10px;
  .graph-map-content{
     width: 100%;
     height:100%;
     background:#fff;
     border-radius:6px;
      .graph-header {
    width:calc(100% - 20px);
    padding:24px 0 14px 30px;
    border-bottom:1px solid #eeeeee;
    .back{
      font-size:20px;
      cursor:pointer;
    }
    .header-text {
      color: #434C6C;
      font-size: 18px;
      font-weight: bold;
      user-select: none;
    }
  }
  .graph-toolbar {
    position: absolute;
    bottom: 20px;
    left: 50%;
    transform: translateX(-50%);
    z-index: 10;
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px;
    background: rgba(255, 255, 255, 0.9);
    border: 1px solid #e4e7ed;
    border-radius: 4px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

    .el-button {
      margin: 0;
    }

    .el-divider {
      margin: 0 4px;
      height: 20px;
    }
  }

  .graph-container {
    width: 100%;
    height: 100%;
  }
  }
}
</style>

