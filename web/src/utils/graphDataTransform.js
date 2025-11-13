function generateColorFromString(str) {
  if (!str) return '#C6E5FF'
  
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash)
  }
  
  const hue = Math.abs(hash % 360)
  const saturation = 60 + (Math.abs(hash) % 20)
  const lightness = 50 + (Math.abs(hash) % 20)
  
  return `hsl(${hue}, ${saturation}%, ${lightness}%)`
}

export function transformGraphData(backendData, options = {}) {
  if (!backendData) {
    return { nodes: [], edges: [] }
  }

  const nodes = backendData.nodes || []
  const edges = backendData.edges || []
  
  const typeColorMap = {}
  nodes.forEach(node => {
    const entityType = node.entity_type || ''
    if (entityType && !typeColorMap[entityType]) {
      typeColorMap[entityType] = generateColorFromString(entityType)
    }
  })
  
  const {
    getNodeId = (node, index) => node.entity_name || `node_${index}`,
    getNodeLabel = (node, index) => node.entity_name || `node_${index}`,
    getNodeSize = (node, index) => node.pagerank ? Math.max(15, Math.min(30, node.pagerank * 100)) : 20,
    getNodeColor = (node, index) => {
      const entityType = node.entity_type || ''
      return typeColorMap[entityType] || '#C6E5FF'
    },
    getEdgeId = (edge, index) => `e${index}`,
    getEdgeLabel = (edge, index) => edge.description || ''
  } = options

  const nodeIdMap = new Map()

  const transformedNodes = nodes.map((node, index) => {
    const nodeId = getNodeId(node, index)
    const nodeLabel = getNodeLabel(node, index)
    const nodeSize = getNodeSize(node, index)
    const nodeColor = getNodeColor(node, index)

    if (node && node.entity_name) {
      nodeIdMap.set(node.entity_name, nodeId)
    }
    if (node && node.entity_id) {
      nodeIdMap.set(String(node.entity_id), nodeId)
    }
    if (node && node.id) {
      nodeIdMap.set(String(node.id), nodeId)
    }

    return {
      ...node,
      id: nodeId,
      label: nodeLabel,
      originalLabel: nodeLabel,
      type: 'circle',
      size: nodeSize,
      style: {
        fill: nodeColor
      }
    }
  })

  const transformedEdges = edges.map((edge, index) => {
    const edgeId = getEdgeId(edge, index)
    const edgeLabel = getEdgeLabel(edge, index)

    const source =
      nodeIdMap.get(edge && edge.source_entity) ||
      nodeIdMap.get(edge && edge.source) ||
      nodeIdMap.get(
        edge && edge.source_id ? String(edge.source_id) : undefined
      ) ||
      (edge && edge.source_entity) ||
      (edge && edge.source) ||
      (edge && edge.source_id ? String(edge.source_id) : undefined) ||
      `source_${index}`

    const target =
      nodeIdMap.get(edge && edge.target_entity) ||
      nodeIdMap.get(edge && edge.target) ||
      nodeIdMap.get(
        edge && edge.target_id ? String(edge.target_id) : undefined
      ) ||
      (edge && edge.target_entity) ||
      (edge && edge.target) ||
      (edge && edge.target_id ? String(edge.target_id) : undefined) ||
      `target_${index}`

    return {
      id: edgeId,
      source,
      target,
      label: edgeLabel,
      ...(edge.weight && {
        style: {
          lineWidth: Math.max(1, Math.min(5, edge.weight / 2))
        }
      }),
      ...edge
    }
  })

  return {
    nodes: transformedNodes,
    edges: transformedEdges
  }
}

export default {
  transformGraphData
}
