class Splinter < Formula
    desc "A cross-platform and language indepedent db migration tool"
    homepage "https://github.com/squareboat/splinter#readme"
    url "https://github.com/squareboat/splinter/archive/refs/tags/0.0.2-alpha.tar.gz"
    sha256 "6f90144fbffbb37370894db2babbf8b8c2274375d1834659344dbb7a72b4b11d"
    license "MIT"
  
    depends_on "go" => :build
  
    def install
      homeDir = Dir::home()
      print "Installing Splinter...\n"
      print "Version: #{version}\n"
      system "go", "build", "-o", "splinter"
      bin.install "splinter"
  
    end
  
    test do
      binFile =File.exist?("#{bin}/splinter")
      assert binFile, "Splinter binary not found"
      output = `#{bin}/splinter --help`
      if output.include? "migration" 
        return true
      end
      return false
    end
  end