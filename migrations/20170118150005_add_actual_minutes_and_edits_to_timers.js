// This migration is for development DB only!
// To run it just paste into shell next script:
// mongo < 20170118150005_add_actual_minutes_and_edits_to_timers.js

conn = new Mongo();
db = conn.getDB("tuna_timer_dev");

DBQuery.shellBatchSize = db.timers.count();
timers = db.timers.find({actual_minutes: null}).toArray();

timers.forEach(function(timer) {
    timer.actual_minutes = NumberInt(timer.minutes);
    timer.edits = [];
    db.timers.save(timer);
});
