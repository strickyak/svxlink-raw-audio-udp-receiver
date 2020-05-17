# svxlink-raw-audio-udp-receiver

## lms.go: A calibration program for Rx volume in svxlink.

Print 28 plus the Log-Mean-Square of groups of 50 samples,
and then ":" and the max of every group of 10 of those.
Listens on UDP for the packet format that comes out of
svxlink when you specify `RAW_AUDIO_UDP_DEST=127.0.0.1:1825`
in your `[Rx1]` clause.  The samples are between -1 and 1,
so their squares are between 0 and 1.  The printed 28+LMS with my
IC-2730 radio with volume at 1:30 and "Mic" device at 42% volume
tend to be
*   7 to 8 when radio squelch is closed,
*   15 to 18 when radio squelch opens but has dead air (receiving empty carrier),
*   19 to 23 when I am speaking,
*   20 to 24.5 for the repeater voice ID,
*   22.7 to 23.2 for the bee-boop,
*   25 is just about the maximum,
*   25.77 on the carrier drop.

Hint: `less -p ':  2' test-and-id-and-drop.42.130.lms`
with 2 spaces between : and 2 to find non-dead air.
