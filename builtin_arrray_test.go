package goja

import "testing"

func TestArrayProtoProp(t *testing.T) {
	const SCRIPT = `
	Object.defineProperty(Array.prototype, '0', {value: 42, configurable: true, writable: false})
	var a = []
	a[0] = 1
	a[0]
	`

	testScript1(SCRIPT, valueInt(42), t)
}

func TestArrayDelete(t *testing.T) {
	const SCRIPT = `
	var a = [1, 2];
	var deleted = delete a[0];
	var undef = a[0] === undefined;
	var len = a.length;

	deleted && undef && len === 2;
	`

	testScript1(SCRIPT, valueTrue, t)
}

func TestArrayDeleteNonexisting(t *testing.T) {
	const SCRIPT = `
	Array.prototype[0] = 42;
	var a = [];
	delete a[0] && a[0] === 42;
	`

	testScript1(SCRIPT, valueTrue, t)
}

func TestArraySetLength(t *testing.T) {
	const SCRIPT = `
	var a = [1, 2];
	var assert0 = a.length == 2;
	a.length = "1";
	a.length = 1.0;
	a.length = 1;
	var assert1 = a.length == 1;
	a.length = 2;
	var assert2 = a.length == 2;
	assert0 && assert1 && assert2 && a[1] === undefined;

	`

	testScript1(SCRIPT, valueTrue, t)
}

func TestArrayReverseNonOptimisable(t *testing.T) {
	const SCRIPT = `
	var a = [];
	Object.defineProperty(a, "0", {get: function() {return 42}, set: function(v) {Object.defineProperty(a, "0", {value: v + 1, writable: true, configurable: true})}, configurable: true})
	a[1] = 43;
	a.reverse();

	a.length === 2 && a[0] === 44 && a[1] === 42;
	`

	testScript1(SCRIPT, valueTrue, t)
}

func TestArrayPushNonOptimisable(t *testing.T) {
	const SCRIPT = `
	Object.defineProperty(Object.prototype, "0", {value: 42});
	var a = [];
	var thrown = false;
	try {
		a.push(1);
	} catch (e) {
		thrown = e instanceof TypeError;
	}
	thrown;
	`

	testScript1(SCRIPT, valueTrue, t)
}

func TestArraySetLengthWithPropItems(t *testing.T) {
	const SCRIPT = `
	var a = [1,2,3,4];
	var thrown = false;

	Object.defineProperty(a, "2", {value: 42, configurable: false, writable: false});
	try {
		Object.defineProperty(a, "length", {value: 0, writable: false});
	} catch (e) {
		thrown = e instanceof TypeError;
	}
	thrown && a.length === 3;
	`

	testScript1(SCRIPT, valueTrue, t)
}

func TestArrayFrom(t *testing.T) {
	const SCRIPT = `
	function checkDestHoles(dest, prefix) {
		assert(dest !== source, prefix + ": dest !== source");
		assert.sameValue(dest.length, 3, prefix + ": dest.length");
		assert.sameValue(dest[0], 1, prefix + ": [0]");
		assert.sameValue(dest[1], undefined, prefix + ": [1]");
		assert(dest.hasOwnProperty("1"), prefix + ': hasOwnProperty("1")');
		assert.sameValue(dest[2], 3, prefix + ": [2]");
	}

	function checkDest(dest, prefix) {
		assert(dest !== source, prefix + ": dest !== source");
		assert.sameValue(dest.length, 3, prefix + ": dest.length");
		assert.sameValue(dest[0], 1, prefix + ": [0]");
		assert.sameValue(dest[1], 2, prefix + ": [1]");
		assert.sameValue(dest[2], 3, prefix + ": [2]");
	}

	var source = [1,2,3];
	var srcHoles = [1,,3];

	checkDest(Array.from(source), "std source/std dest");
	checkDestHoles(Array.from(srcHoles), "std source (holes)/std dest");

	function Iter() {
		this.idx = 0;
	}
	Iter.prototype.next = function() {
		if (this.idx < source.length) {
			return {value: source[this.idx++]};
		} else {
			return {done: true};
		}
	}

	var src = {};
	src[Symbol.iterator] = function() {
		return new Iter();
	}
	checkDest(Array.from(src), "iter src/std dest");

	src = {0: 1, 2: 3, length: 3};
	checkDestHoles(Array.from(src), "arrayLike src/std dest");

	function A() {}
	A.from = Array.from;

	checkDest(A.from(source), "std src/cust dest");
	checkDestHoles(A.from(srcHoles), "std src (holes)/cust dest");
	checkDestHoles(A.from(src), "arrayLike src/cust dest");

	function T2() {
	  Object.defineProperty(this, 0, {
		configurable: false,
		writable: true,
		enumerable: true
	  });
	}

	assert.throws(TypeError, function() {
		Array.from.call(T2, source);
	});

	`

	testScript1(TESTLIB+SCRIPT, _undefined, t)
}

func TestArrayOf(t *testing.T) {
	const SCRIPT = `
	function T1() {
	  Object.preventExtensions(this);
	}
	
	assert.throws(TypeError, function() {
	  Array.of.call(T1, 'Bob');
	});

	function T2() {
	  Object.defineProperty(this, 0, {
		configurable: false,
		writable: true,
		enumerable: true
	  });
	}
	
	assert.throws(TypeError, function() {
	  Array.of.call(T2, 'Bob');
	})

	result = Array.of.call(undefined);
	assert(
	  result instanceof Array,
	  'this is not a constructor'
	);

	result = Array.of.call(Math.cos);
	assert(
	  result instanceof Array,
	  'this is a builtin function with no [[Construct]] slot'
	);

	`

	testScript1(TESTLIB+SCRIPT, _undefined, t)
}
