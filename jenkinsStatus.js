const request = require('request');

const jenkinsUrl = process.argv[2];

function handleError(error) {
    console.error(error);
    console.log('Error');
}

const goodResults = 'SUCCESS'

request(jenkinsUrl, (error, response, rawBody) => {
    if (!error && response.statusCode >= 200 && response.statusCode < 300) {
        const body = JSON.parse(rawBody);
        const building = ('building' in body && body.building);

        let result = 'PENDING';
        if (!building) {
            if ('result' in body && body.result != null) {
                result = body.result;
            }
            else {
                result = 'UNKNOWN';
            }
        }

        const prettyResult = result.replace('_', ' ').toLowerCase().replace(/\b\w/g, l => l.toUpperCase());

        console.log(prettyResult);
    }
    else {
        handleError(error);
    }
});