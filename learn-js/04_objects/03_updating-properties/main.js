function addID(campaignRecord) {
  campaignRecord.id = `${campaignRecord.campaignName}-${campaignRecord.senderName}`
  return campaignRecord;
}

// don't touch below this line

export { addID };

