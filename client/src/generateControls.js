export default function generateControls() {
    let controls = {
      0: {
        i: 0,
        p: null,
        t: 'Page',
        counter: 0,
        expanded: true,
        c: []
      }
    }
  
    for (let i = 1; i < 100; i++) {
      let parentId = Math.floor(Math.pow(Math.random(), 2) * i)
      controls[i] = {
        i: i,
        p: parentId,
        t: 'Node',
        counter: 0,
        expanded: true,
        c: []
      }
      controls[parentId].c.push(i)
    }
  
    return controls
  }
  