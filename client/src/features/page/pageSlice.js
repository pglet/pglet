import { createSlice } from '@reduxjs/toolkit'

let nextId = 0

const initialState = {
    "name": "test-1",
    "controls": {
        "0": {
            "c": [
                "1"
            ],
            "i": "0",
            "p": "",
            "t": "Page"
        },
        "1": {
            "c": [
                "2",
                "3"
            ],
            "i": "1",
            "p": "0",
            "t": "Row"
        },
        "2": {
            "c": [
                "myTxt",
                "myBtn"
            ],
            "i": "2",
            "p": "1",
            "t": "Column"
        },
        "3": {
            "c": [
                "5"
            ],
            "i": "3",
            "p": "1",
            "t": "Column"
        },
        "myTxt": {
            "i": "myTxt",
            "p": "2",
            "t": "Text",
            "text": "Hello, world!"
        },
        "5": {
            "i": "5",
            "p": "3",
            "t": "Button",
            "text": "Click me!"
        },
        "myBtn": {
            "i": "myBtn",
            "p": "2",
            "t": "Button",
            "text": "Cancel"
        }
    }
}

const pageSlice = createSlice({
    name: 'page',
    initialState: initialState,
    reducers: {
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
            const { nodeId, newProps } = action.payload
            const node = state.controls[nodeId]
            Object.assign(node, newProps)
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

export const { createNode, increment, toggleExpand, addChild, removeChild, changeProps, deleteNode } = pageSlice.actions

export default pageSlice.reducer