const hasHighOpenRate = true;
const isRecent = true;
const hasStrongReplyRate = false;
const canBeResent = true;
const isFlaggedAsSpam = false;

// don't touch above this line

const isHighEngagement = hasHighOpenRate && isRecent && (hasStrongReplyRate || canBeResent) && !isFlaggedAsSpam;
// don't touch below this line

console.log(`The campaign is high-engagement: ${isHighEngagement}`);

