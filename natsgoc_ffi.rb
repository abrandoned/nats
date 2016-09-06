require 'rubygems'
require 'ffi'
require 'thread'

module NATSFFI
  extend FFI::Library
  ffi_lib File.expand_path("./libnatsc.so", File.dirname(__FILE__))

  attach_function :Connect, [:string], :int32
  attach_function :Close, [:int], :void
  attach_function :CloseAll, [], :void
  attach_function :Flush, [:int], :void
  attach_function :FlushAll, [], :void
  attach_function :Publish, [:int, :string, :string, :int], :void
  attach_function :Request, [:int, :string, :string], :string

  def self.test_threaded
    start = Time.now
    num_threads = 4
    publish_per_thread = 2_000_000
    threads = []
    subject = "hello"
    message = "world"
    message_size = message.size

    num_threads.times do
      threads << Thread.new do
        handle = NATSFFI.Connect("nats://localhost:4222")

        publish_per_thread.times do
          NATSFFI.Publish(handle, subject, message, message_size)
        end

        NATSFFI.Flush(handle)
        NATSFFI.Close(handle)
      end
    end

    threads.map(&:join)
    finish = Time.now
    puts <<-FINISH
    THREADS: #{num_threads}
    PUBLISH PER THREAD: #{publish_per_thread}
    START: #{start}
    FINISH: #{finish}
    PER SECOND: #{(num_threads * publish_per_thread)/(finish.to_i - start.to_i)}
    FINISH
  end
end
