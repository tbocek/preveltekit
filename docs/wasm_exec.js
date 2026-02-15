// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This file has been modified for use by the TinyGo compiler.
// Browser-only version.

(() => {
  const global = window;
  const encoder = new TextEncoder("utf-8");
  const decoder = new TextDecoder("utf-8");

  const enosys = () => {
    const err = new Error("not implemented");
    err.code = "ENOSYS";
    return err;
  };

  global.fs = new Proxy(
    {
      constants: {
        O_WRONLY: -1,
        O_RDWR: -1,
        O_CREAT: -1,
        O_TRUNC: -1,
        O_APPEND: -1,
        O_EXCL: -1,
      },
    },
    {
      get(target, prop) {
        if (prop in target) return target[prop];
        return (...args) => {
          const callback = args[args.length - 1];
          if (typeof callback === "function") callback(enosys());
          else throw enosys();
        };
      },
    },
  );

  global.process = new Proxy(
    {
      pid: -1,
      ppid: -1,
    },
    {
      get(target, prop) {
        if (prop in target) return target[prop];
        return () => {
          throw enosys();
        };
      },
    },
  );

  let reinterpretBuf = new DataView(new ArrayBuffer(8));
  var logLine = [];
  const wasmExit = {};

  global.Go = class {
    constructor() {
      const mem = () => {
        return new DataView(this._inst.exports.memory.buffer);
      };

      const unboxValue = (v_ref) => {
        reinterpretBuf.setBigInt64(0, v_ref, true);
        const f = reinterpretBuf.getFloat64(0, true);
        if (f === 0) {
          return undefined;
        }
        if (!isNaN(f)) {
          return f;
        }
        const id = v_ref & 0xffffffffn;
        return this._values[id];
      };

      const loadValue = (addr) => {
        return unboxValue(mem().getBigUint64(addr, true));
      };

      const boxValue = (v) => {
        const nanHead = 0x7ff80000n;

        if (typeof v === "number") {
          if (isNaN(v)) {
            return nanHead << 32n;
          }
          if (v === 0) {
            return (nanHead << 32n) | 1n;
          }
          reinterpretBuf.setFloat64(0, v, true);
          return reinterpretBuf.getBigInt64(0, true);
        }

        switch (v) {
          case undefined:
            return 0n;
          case null:
            return (nanHead << 32n) | 2n;
          case true:
            return (nanHead << 32n) | 3n;
          case false:
            return (nanHead << 32n) | 4n;
        }

        let id = this._ids.get(v);
        if (id === undefined) {
          id = this._idPool.pop();
          if (id === undefined) {
            id = BigInt(this._values.length);
          }
          this._values[id] = v;
          this._goRefCounts[id] = 0;
          this._ids.set(v, id);
        }
        this._goRefCounts[id]++;
        let typeFlag = 1n;
        switch (typeof v) {
          case "string":
            typeFlag = 2n;
            break;
          case "symbol":
            typeFlag = 3n;
            break;
          case "function":
            typeFlag = 4n;
            break;
        }
        return id | ((nanHead | typeFlag) << 32n);
      };

      const storeValue = (addr, v) => {
        let v_ref = boxValue(v);
        mem().setBigUint64(addr, v_ref, true);
      };

      const loadSlice = (array, len, cap) => {
        return new Uint8Array(this._inst.exports.memory.buffer, array, len);
      };

      const loadSliceOfValues = (array, len, cap) => {
        const a = new Array(len);
        for (let i = 0; i < len; i++) {
          a[i] = loadValue(array + i * 8);
        }
        return a;
      };

      const loadString = (ptr, len) => {
        return decoder.decode(
          new DataView(this._inst.exports.memory.buffer, ptr, len),
        );
      };

      const timeOrigin = Date.now() - performance.now();
      this.importObject = {
        wasi_snapshot_preview1: {
          fd_write: function (fd, iovs_ptr, iovs_len, nwritten_ptr) {
            let nwritten = 0;
            if (fd == 1) {
              for (let iovs_i = 0; iovs_i < iovs_len; iovs_i++) {
                let iov_ptr = iovs_ptr + iovs_i * 8;
                let ptr = mem().getUint32(iov_ptr + 0, true);
                let len = mem().getUint32(iov_ptr + 4, true);
                nwritten += len;
                for (let i = 0; i < len; i++) {
                  let c = mem().getUint8(ptr + i);
                  if (c == 13) {
                  } else if (c == 10) {
                    let line = decoder.decode(new Uint8Array(logLine));
                    logLine = [];
                    console.log(line);
                  } else {
                    logLine.push(c);
                  }
                }
              }
            } else {
              console.error("invalid file descriptor:", fd);
            }
            mem().setUint32(nwritten_ptr, nwritten, true);
            return 0;
          },
          fd_close: () => 0,
          fd_fdstat_get: () => 0,
          fd_seek: () => 0,
          proc_exit: (code) => {
            this.exited = true;
            this.exitCode = code;
            this._resolveExitPromise();
            throw wasmExit;
          },
          random_get: (bufPtr, bufLen) => {
            crypto.getRandomValues(loadSlice(bufPtr, bufLen));
            return 0;
          },
        },
        gojs: {


          "syscall/js.finalizeRef": (v_ref) => {
            const id = v_ref & 0xffffffffn;
            if (this._goRefCounts?.[id] !== undefined) {
              this._goRefCounts[id]--;
              if (this._goRefCounts[id] === 0) {
                const v = this._values[id];
                this._values[id] = null;
                this._ids.delete(v);
                this._idPool.push(id);
              }
            }
          },
          // end

          "syscall/js.stringVal": (value_ptr, value_len) => {
            value_ptr >>>= 0;
            const s = loadString(value_ptr, value_len);
            return boxValue(s);
          },
          // end

          "syscall/js.valueGet": (v_ref, p_ptr, p_len) => {
            let prop = loadString(p_ptr, p_len);
            let v = unboxValue(v_ref);
            let result = Reflect.get(v, prop);
            return boxValue(result);
          },
          // end

          "syscall/js.valueSet": (v_ref, p_ptr, p_len, x_ref) => {
            const v = unboxValue(v_ref);
            const p = loadString(p_ptr, p_len);
            const x = unboxValue(x_ref);
            Reflect.set(v, p, x);
          },
          // end


          "syscall/js.valueIndex": (v_ref, i) => {
            return boxValue(Reflect.get(unboxValue(v_ref), i));
          },
          // end

          "syscall/js.valueSetIndex": (v_ref, i, x_ref) => {
            Reflect.set(unboxValue(v_ref), i, unboxValue(x_ref));
          },
          // end

          "syscall/js.valueCall": (
            ret_addr,
            v_ref,
            m_ptr,
            m_len,
            args_ptr,
            args_len,
            args_cap,
          ) => {
            const v = unboxValue(v_ref);
            const name = loadString(m_ptr, m_len);
            const args = loadSliceOfValues(args_ptr, args_len, args_cap);
            try {
              const m = Reflect.get(v, name);
              storeValue(ret_addr, Reflect.apply(m, v, args));
              mem().setUint8(ret_addr + 8, 1);
            } catch (err) {
              storeValue(ret_addr, err);
              mem().setUint8(ret_addr + 8, 0);
            }
          },
          // end


          "syscall/js.valueNew": (
            ret_addr,
            v_ref,
            args_ptr,
            args_len,
            args_cap,
          ) => {
            const v = unboxValue(v_ref);
            const args = loadSliceOfValues(args_ptr, args_len, args_cap);
            try {
              storeValue(ret_addr, Reflect.construct(v, args));
              mem().setUint8(ret_addr + 8, 1);
            } catch (err) {
              storeValue(ret_addr, err);
              mem().setUint8(ret_addr + 8, 0);
            }
          },
          // end

          "syscall/js.valueLength": (v_ref) => {
            return unboxValue(v_ref).length;
          },
          // end

          "syscall/js.valuePrepareString": (ret_addr, v_ref) => {
            const s = String(unboxValue(v_ref));
            const str = encoder.encode(s);
            storeValue(ret_addr, str);
            mem().setInt32(ret_addr + 8, str.length, true);
          },
          // end

          "syscall/js.valueLoadString": (
            v_ref,
            slice_ptr,
            slice_len,
            slice_cap,
          ) => {
            const str = unboxValue(v_ref);
            loadSlice(slice_ptr, slice_len, slice_cap).set(str);
          },
          // end



        },
      };

      this.importObject.env = this.importObject.gojs;
    }

    async run(instance) {
      this._inst = instance;
      this._values = [NaN, 0, null, true, false, global, this];
      this._goRefCounts = [];
      this._ids = new Map();
      this._idPool = [];
      this.exited = false;
      this.exitCode = 0;

      if (this._inst.exports._start) {
        let exitPromise = new Promise((resolve, reject) => {
          this._resolveExitPromise = resolve;
        });

        try {
          this._inst.exports._start();
        } catch (e) {
          if (e !== wasmExit) throw e;
        }

        await exitPromise;
        return this.exitCode;
      } else {
        this._inst.exports._initialize();
      }
    }

    _resume() {
      if (this.exited) {
        throw new Error("Go program has already exited");
      }
      try {
        this._inst.exports.resume();
      } catch (e) {
        if (e !== wasmExit) throw e;
      }
      if (this.exited) {
        this._resolveExitPromise();
      }
    }

    _makeFuncWrapper(id) {
      const go = this;
      return function () {
        const event = { id: id, this: this, args: arguments };
        go._pendingEvent = event;
        go._resume();
        return event.result;
      };
    }
  };
})();
