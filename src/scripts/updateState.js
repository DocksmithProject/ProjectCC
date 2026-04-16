const fs = require('fs');

let data = JSON.parse(fs.readFileSync('./state/index.json', 'utf-8'));

data.currentState = "processing";
data.version++;

fs.writeFileSync('./state/index.json', JSON.stringify(data, null, 2));

console.log("State updated");
