set -e

mongo <<EOF
use systems-management

db.createCollection("user", { capped: false });
db.createCollection("domain", { capped: false });
db.user.createIndex({ email: 1 }, { unique: true });
db.user.insert([
  {
    email: "$MONGO_SUPERADMIN_EMAIL",
    role: "superadmin",
    created: parseInt(new Date().getTime() / 1000),
    password: hex_md5("$MONGO_SUPERADMIN_PASSWORD"),
  },
]);
EOF
