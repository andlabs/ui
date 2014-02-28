// 27 february 2014
#include <stdio.h>
#include <objc/message.h>
#include <objc/objc.h>
#include <objc/runtime.h>

int main(void)
{
	id NSString = objc_getClass("NSString");
	SEL stringFromUTF8String =
		sel_getUid("stringWithUTF8String:");
	id str = objc_msgSend(NSString,
		stringFromUTF8String,
		"hello, world\n");
	SEL UTF8String =
		sel_getUid("UTF8String");

	printf("%s",
		(char *) objc_msgSend(str,
			UTF8String));
	return 0;
}
