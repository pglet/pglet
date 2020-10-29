import { createSlice, current } from '@reduxjs/toolkit'

let nextId = 0

const initialState = {
    "name": "test-1",
    "error": null,
    "controls": {
        "page": {
            "c": [],
            "i": "page",
            "p": "",
            "t": "page"
        }
    }
}

const pageSlice = createSlice({
    name: 'page',
    initialState: initialState,
    reducers: {
        registerWebClientSuccess(state, action) {
            state.loading = false;
            state.sessionId = action.payload.id;
            state.controls = action.payload.controls;
        },
        registerWebClientError(state, action) {
            state.loading = false;
            state.error = action.payload;
        },
        addPageControlsSuccess(state, action) {
            action.payload.forEach(ctrl => {
                if (!state.controls[ctrl.i]) {
                    state.controls[ctrl.i] = ctrl;
                    state.controls[ctrl.p].c.push(ctrl.i)
                }
            })
        },
        addPageControlsError(state, action) {
            state.error = action.payload;
        },        
        createNode: {
            reducer(state, action) {
                const { nodeId, parentId } = action.payload
                state.controls[nodeId] = {
                    i: nodeId,
                    p: parentId,
                    t: 'Node',
                    c: [],
                    counter: 0,
                    expanded: true
                }
                return nodeId;
            },
            prepare(parentId) {
                return {
                    payload: {
                        parentId,
                        nodeId: `new_${nextId++}`
                    }
                }
            }
        },
        increment(state, action) {
            const node = state.controls[action.payload]
            node.counter++
        },
        toggleExpand(state, action) {
            const node = state.controls[action.payload]
            node.expanded = !node.expanded
        },
        addChild(state, action) {
            const { nodeId, childId } = action.payload
            const node = state.controls[nodeId]
            node.c.push(childId)
        },
        removeChild(state, action) {
            const { nodeId, childId } = action.payload
            const node = state.controls[nodeId]
            node.c = node.c.filter(id => id !== childId)
        },
        changeProps(state, action) {

            action.payload.forEach(props => {
                const ctrl = state.controls[props.i];
                if (ctrl) {
                    Object.assign(ctrl, props)
                }
            })
            //console.log(current(state))
        },
        deleteNode(state, action) {
            const nodeId = action.payload
            const descendantIds = getAllDescendantIds(state.controls, nodeId)
            return deleteMany(state.controls, [nodeId, ...descendantIds])
        }
    }
})

const getAllDescendantIds = (controls, nodeId) => {
    if (controls[nodeId].c) {
        return controls[nodeId].c.reduce((acc, childId) => (
            [...acc, childId, ...getAllDescendantIds(controls, childId)]
        ), [])
    }
    return []
}

const deleteMany = (controls, ids) => {
    ids.forEach(id => delete controls[id])
}

export const {
    registerWebClientSuccess,
    registerWebClientError,
    addPageControlsSuccess,
    addPageControlsError,
    createNode,
    increment,
    toggleExpand,
    addChild,
    removeChild,
    changeProps,
    deleteNode
} = pageSlice.actions

export default pageSlice.reducer