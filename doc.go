/*Package attache: an arsenal for the web.

Attachè (mispronounced for fun as "Attachey", said like "Apache")
is a new kind of web framework for go. It allows you to define your applications
in terms of methods on a context type that is instantiated for each request.

An Attachè Application is bootstrapped from a context type that implements the Context interface.
In order to enforce that the context type is a struct type at compile time,
embedding the BaseContext type is necessary to satisfy this interface.
Functionality is provided by additional methods implemented on the context type.

Attachè provides several capabilities (optional interfaces) that a Context type can implement.
These interfaces take the form of "Has...". They are documented individually.

Aside from these capabilites, a Context can provide several other types of methods:

1. Route methods
Routing methods must follow the naming convention of <HTTP METHOD>_<[PathInCamelCase]>
PathInCamelCase is optional, and defaults to the empty string.
In this case, the method name would still need to include the underscore (i.e. GET_).
The path which a routing method serves is determined by converting the PathInCamelCase
to snake case, and then replacing underscores with the path separator.
Examples:
	(empty)				-> matches /
	Index				-> matches /index
	TwoWords			-> matches /two/words
	Ignore_underscores	-> matches /ignoreunderscores (for making long names easier to read without affecting the path)
	TESTInitialism		-> matches /test/initialism

Route methods can take any number/types of arguments. Values are provided by the Application's
dependency injection context for this request. If a value isn't available for an argument's type,
the zero-value is used. Return values are unchecked, and really shouldn't be provided.
That limitation is not enforced in the code, to provide flexibility.

If a Context type provides a BEFORE_ method matching the name of the Route method,
it will be run immediately before the Route method in the handler stack

If a Context type provides an AFTER_ method matching the name of the Route method,
it will be run immediately after the Route method in the handler stack. Although there is no
hard enforcement, an AFTER_ method should assume the ResponseWriter has been closed.

2. Mount methods
Mount methods must follow the naming convention of MOUNT_<[PathInCamelCase]>
and have the signature func() (http.Handler, error).
They are only called once, on an uninitialized Context during bootstrapping.
If an error is returned (or the method panics), bootstrapping will fail. If there is no error,
the returned http.Handler is mounted at the path specified by PathInCamelCase, using the same method as
Routing methods. The mounted handler is called with the mounted prefix stripped from the request.

3. Guard methods
Mount methods must follow the naming convention of GUARD_<[PathInCamelCase]>.
The path at which a Guard method is applied is determined from PathInCamelCase, in the same way as a Route method.
A Guard method will run for the path at which it is registered, as well as any paths that contain the guarded path.
To prevent the handler stack from continuing execution, the Guard method must either panic or call one of the
Attachè helper methods (Error, ErrorFatal, ErrorMessage, etc...). If there are multiple Guard methods on a path,
they are run in they order they are encountered (i.e. the guards for shorter paths run before the guards for longer paths).

4. Provider methods
Provider methods must follow the naming convention of PROVIDE_<UniqueDescriptiveName>
and have the signature func(*http.Request) interface{}.
Provider methods are called when a route or a guarded mount match the request's path.
The returned value is available for injection into any handlers in the handler stack
for the current request. Because of the frequency with which these are called, it is best
to define as few as possible.

Attachè also comes with a CLI.
Currently, the supported commands are `attache new` and `attache gen`. Please see README.md for more info.
*/
package attache
