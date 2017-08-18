const os = require('os-utils');

os.cpuUsage((usage) => {
    const percentUsage = Math.round(usage * 100) + '';

    console.log(percentUsage.padStart(2, '0') + '%');
});