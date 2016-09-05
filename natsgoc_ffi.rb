require 'rubygems'
require 'ffi'
require 'thread'

module NATSFFI
  extend FFI::Library
  ffi_lib File.expand_path("./libnatsc.so", File.dirname(__FILE__))

  attach_function :Connect, [:string], :long_long
  attach_function :Close, [:long_long], :void
  attach_function :CloseAll, [], :void
  attach_function :Flush, [:long_long], :void
  attach_function :FlushAll, [], :void
  attach_function :Publish, [:long_long, :string, :string], :void
  attach_function :Request, [:long_long, :string, :string], :string 

  def self.test_threaded
    threads = []

    8.times do
      threads << Thread.new do
        conn = NATSFFI.Connect("nats://localhost:4222")

        1_000_000.times do
          NATSFFI.Publish(conn, "hello", "world")
        end

        NATSFFI.Close(conn)
      end
    end

    threads.map(&:join)
  end
end
