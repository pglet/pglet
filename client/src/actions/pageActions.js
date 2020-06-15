export const INCREMENT = 'INCREMENT'
export const TOGGLE_EXPAND = 'TOGGLE_EXPAND'
export const CREATE_NODE = 'CREATE_NODE'
export const DELETE_NODE = 'DELETE_NODE'
export const ADD_CHILD = 'ADD_CHILD'
export const REMOVE_CHILD = 'REMOVE_CHILD'
export const CHANGE_PROPS = 'CHANGE_PROPS'

export const increment = (nodeId) => ({
  type: INCREMENT,
  nodeId
})

export const toggleExpand = (nodeId) => ({
  type: TOGGLE_EXPAND,
  nodeId
})

let nextId = 0
export const createNode = (parentId) => ({
  type: CREATE_NODE,
  parentId,
  nodeId: `new_${nextId++}`
})

export const deleteNode = (nodeId) => ({
  type: DELETE_NODE,
  nodeId
})

export const addChild = (nodeId, childId) => ({
  type: ADD_CHILD,
  nodeId,
  childId
})

export const removeChild = (nodeId, childId) => ({
  type: REMOVE_CHILD,
  nodeId,
  childId
})

export const changeProps = (nodeId, newProps) => ({
  type: CHANGE_PROPS,
  nodeId,
  newProps
})
