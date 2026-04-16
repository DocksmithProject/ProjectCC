const fs = require('fs');

const data = fs.readFileSync('./state/index.json', 'utf-8');
console.log(JSON.parse(data));
