const fs = require('fs');

fs.writeFileSync('./state/index.json', JSON.stringify({
  currentState: "initial",
  version: 1
}, null, 2));

console.log("State reset done");
