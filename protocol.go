package adb

/*
The smart protocol is documented in SERVICES.TXT.
https://android.googlesource.com/platform/packages/modules/adb/+/HEAD/SERVICES.TXT
*/

// Request make to the ADB server
func (c *Conn) Request(req string) error) {
	if err := c.SendMessage(req); err != nil {
		return c, err
	}
	return c, nil
}

// HOST SERVICES:

// Version host:version
// Ask the ADB server for its internal version number.
func (c *Conn) Version() string {
	req := "host:version"
	_, c.Request(req)
}

// Host
/*
switch the transport to a real device,
request without a prefix "host" will add automatically.


*/
func (c *Conn) Host(request string) {
}

// LOCAL SERVICES:
