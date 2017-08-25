const request = require('request');
const emoji = require('node-emoji');

// this uses this jenkins plugin https://github.com/jenkinsci/pipeline-stage-view-plugin

// the expected url format is https://<jenkins-host>/job/<job-name>/lastBuild/wfapi/describe
const jenkinsUrl = process.argv[2];
const emojiName = (process.argv.length > 2) ? process.argv[3] : null;

const emojiPrefix = (emojiName != null) ? ((emoji.hasEmoji(emojiName) ? emoji.get(emojiName) : emoji.get('x'))) : '';

function handleError(error) {
    console.error(error);
    console.log('Error');
}

const StatusMap = {
    NOT_EXECUTED: emoji.get('arrows_counterclockwise') + 'Waiting',
    ABORTED: emoji.get('heavy_multiplication_x') + 'Aborted',
    SUCCESS: emoji.get('vertical_traffic_light'),
    IN_PROGRESS: emoji.get('arrows_counterclockwise') + 'Running',
    PAUSED_PENDING_INPUT: emoji.get('double_vertical_bar') + 'Paused',
    FAILED: emoji.get('red_circle') + 'Failed',
    UNSTABLE: emoji.get('question') + 'Unstable'
}

request(jenkinsUrl, (error, response, rawBody) => {
    if (!error && response.statusCode >= 200 && response.statusCode < 300) {
        const body = JSON.parse(rawBody);
        const status = ('status' in body) ? body.status : 'UNKNOWN';
        const statusMessage = (status in StatusMap) ? StatusMap[status] : emoji.get('question') + 'Unknown';

        console.log(emojiPrefix + statusMessage);
    }
    else {
        handleError(error);
    }
});