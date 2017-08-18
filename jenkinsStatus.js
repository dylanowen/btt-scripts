const request = require('request');
const emoji = require('node-emoji');

const jenkinsUrl = process.argv[2];
const showEmoji = process.argv.length > 2 && process.argv[3] == 'true';

function handleError(error) {
    console.error(error);
    console.log('Error');
}

const ResultMap = {
    success: emoji.get('white_check_mark'),
    aborted: emoji.get('black_circle'),
    not_built: emoji.get('warning'),
    unstable: emoji.get('warning'),
    failure: emoji.get('red_circle'),
    pending: emoji.get('large_blue_circle'),
    unknown: emoji.get('black_circle')
}

request(jenkinsUrl, (error, response, rawBody) => {
    if (!error && response.statusCode >= 200 && response.statusCode < 300) {
        const body = JSON.parse(rawBody);
        const building = ('building' in body && body.building);

        let result = 'pending';
        if (!building) {
            if ('result' in body && body.result != null) {
                result = body.result.toLowerCase();
            }
            else {
                result = 'unknown';
            }
        }

        let prettyResult;
        if (showEmoji) {
            prettyResult = (result in ResultMap) ? ResultMap[result] : ResultMap.unknown;
        }
        else {
            prettyResult = result.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase());
        }

        console.log(prettyResult);
    }
    else {
        handleError(error);
    }
});