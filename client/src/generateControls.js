export default function generateControls() {
    let controls = {
      0: {
        id: 0,
        parentId: null,
        type: 'Page',
        counter: 0,
        expanded: true,
        childIds: []
      }
    }
  
    for (let i = 1; i < 100; i++) {
      let parentId = Math.floor(Math.pow(Math.random(), 2) * i)
      controls[i] = {
        id: i,
        parentId: parentId,
        type: 'Node',
        counter: 0,
        expanded: true,
        childIds: []
      }
      controls[parentId].childIds.push(i)
    }
  
    return controls
  }
  