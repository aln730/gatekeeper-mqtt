import { Router } from "express";
const router = Router();

router.get("/", async (req, res) => {
  const { cursor, since, until } = req.query;
  const query = {};
  if (since || until || cursor) {
    query.timestamp = {};
    if (cursor) query.timestamp.$lt = new Date(cursor);
    if (since)  query.timestamp.$gte = new Date(since);
    if (until)  query.timestamp.$lte = new Date(until);
  }

  const logs = await req.ctx.db
    .collection("accessLogs")
    .find(query)
    .sort({ timestamp: -1 })
    .limit(50)
    .toArray();

  const nextCursor = logs.length === 50
    ? logs[logs.length - 1].timestamp.toISOString()
    : null;

  res.json({ logs, cursor: nextCursor });
});

export default router;