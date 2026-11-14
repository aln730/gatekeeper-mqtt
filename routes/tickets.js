import { Router } from "express";
import { ObjectId } from "mongodb";

const router = Router();

const TYPES = {
    groups: { collection: "groupTickets", subject: "groupId" },
    users: { collection: "userTickets", subject: "userId" },
};

router.param("type", (req, res, next, name) => {
    if (!Object.hasOwn(TYPES, name)) {
        return res.status(404).json({ message: "type not found" });
    }
    req.type = TYPES[name];
    next();
});

router.param("id", (req, res, next, id) => {
    if (!ObjectId.isValid(id)) {
        return res.status(422).json({ message: "invalid ticket id" });
    }
    req.ticketId = new ObjectId(id);
    next();
});

function parseTicket(req, res, requireAll) {
    const body = req.body ?? {};
    const checks = {
        doorId: (v) => typeof v === "string" && v !== "",
        [req.type.subject]: (v) => typeof v === "string" && v !== "",
        priority: Number.isInteger,
        granted: (v) => typeof v === "boolean",
    };

    const fields = {};
    for (const [key, isValid] of Object.entries(checks)) {
        if (!(key in body)) {
            if (requireAll) {
                res.status(422).json({
                    message: `No '${key}' field specified`,
                });
                return null;
            }
            continue;
        }
        if (!isValid(body[key])) {
            res.status(422).json({ message: `Invalid '${key}' field` });
            return null;
        }
        fields[key] = body[key];
    }
    if (Object.keys(fields).length === 0) {
        res.status(422).json({ message: "No fields specified" });
        return null;
    }
    return fields;
}

router.get("/:type", async (req, res) => {
    const tickets = await req.ctx.db
        .collection(req.type.collection)
        .find({})
        .sort({ priority: -1 })
        .toArray();
    res.json(tickets);
});

router.put("/:type/:id", async (req, res) => {
    const fields = parseTicket(req, res, true);
    if (!fields) return;

    await req.ctx.db
        .collection(req.type.collection)
        .replaceOne({ _id: req.ticketId }, fields, { upsert: true });
    res.status(204).send(null);
});

router.patch("/:type/:id", async (req, res) => {
    const fields = parseTicket(req, res, false);
    if (!fields) return;

    const result = await req.ctx.db
        .collection(req.type.collection)
        .updateOne({ _id: req.ticketId }, { $set: fields });

    if (result.matchedCount === 0) {
        return res.status(404).json({ message: "ticket not found" });
    }

    res.status(204).send(null);
});

router.delete("/:type/:id", async (req, res) => {
    const result = await req.ctx.db
        .collection(req.type.collection)
        .deleteOne({ _id: req.ticketId });
    if (result.deletedCount === 0) {
        return res.status(404).json({ message: "ticket not found" });
    }

    res.status(204).send(null);
});

export default router;
