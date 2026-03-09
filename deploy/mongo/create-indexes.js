// MongoDB index setup for Tarot Daily MVP.
// Run from mongosh with:
//   mongosh "<MONGO_URI>" --file deploy/mongo/create-indexes.js
// Optional DB override:
//   mongosh "<MONGO_URI>" --eval "const DB_NAME='tarot_daily'" --file deploy/mongo/create-indexes.js

const targetDbName = typeof DB_NAME !== "undefined" && DB_NAME ? DB_NAME : "tarot_daily";
const targetDb = db.getSiblingDB(targetDbName);

print(`Using database: ${targetDbName}`);

// users.login unique
targetDb.users.createIndex(
  { login: 1 },
  {
    name: "uniq_users_login",
    unique: true
  }
);

// daily_readings(userId, date) unique
targetDb.daily_readings.createIndex(
  { userId: 1, date: 1 },
  {
    name: "uniq_daily_readings_user_date",
    unique: true
  }
);

// cards.code unique (contract-level requirement from data model)
targetDb.cards.createIndex(
  { code: 1 },
  {
    name: "uniq_cards_code",
    unique: true
  }
);

print("Indexes created/ensured:");
printjson(targetDb.users.getIndexes());
printjson(targetDb.daily_readings.getIndexes());
printjson(targetDb.cards.getIndexes());
