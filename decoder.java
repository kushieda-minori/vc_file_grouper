/**
 * @Author Kellindil Maendellyn
 * https://valkyriecrusade.fandom.com/wiki/Thread:119497#19
 */
package vc;

import java.io.IOException;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.StandardOpenOption;

public class Decoder {
	public static void main(String[] args) throws IOException {
		Path imageDir = Paths.get("hd");
		try(DirectoryStream<Path> files = Files.newDirectoryStream(imageDir)) {
			for (Path child : files) {
				if (child.toFile().isFile()) {
					decodeDataFile(child.toString());
				}
			}
		}
	}

	// File header : 16 bytes
	// 4 bytes for the signature (CODE)
	// 8 bytes of unknown data
	// 4 bytes for one of the encoding's keys (the second key is a magic number
	// known from the app, 0x45AF6E5D at the time of writing)
	//
	// The remainder of the file is encoded 4 bytes by 4 bytes, the last few
	// bytes unencoded if the file's length is not a multiple of 4
	private static byte[] decodeDataFile(String filePath) throws IOException {
		Path path = Paths.get(filePath);
		byte[] bytes = Files.readAllBytes(path);
		// File signature is always "CODE", or file is not encoded
		if ((bytes[0] & 0xFF) == 'C' && (bytes[1] & 0xFF) == 'O'
				&& (bytes[2] & 0xFF) == 'D' && (bytes[3] & 0xFF) == 'E') {
			int subMe = toInt(bytes, 12);
			int xorMe = 0x45AF6E5D;

			// We'll ignore the 16-bytes signature
			int excessBytes = bytes.length % 4;
			int encodedLength = bytes.length - 16 - excessBytes;
			ByteBuffer result = ByteBuffer
					.allocate(encodedLength + excessBytes);

			// decode 4 by 4
			for (int i = 0; i < encodedLength / 4; i++) {
				int decodedBytes = (toInt(bytes, 16 + (i * 4)) ^ xorMe) - subMe;

				ByteBuffer buf = ByteBuffer.allocate(4);
				buf.order(ByteOrder.LITTLE_ENDIAN);
				buf.putInt(decodedBytes);
				buf.flip();
				result.put(buf.array());
			}

			// copy the last few bytes as-is
			if (excessBytes > 0) {
				byte[] remainder = new byte[excessBytes];
				System.arraycopy(bytes, 16 + encodedLength, remainder, 0,
						excessBytes);

				result.put(remainder);
			}

			result.flip();
			byte[] resultArray = result.array();

			final Path output;
			if ((resultArray[0] & 0xFF) == 0x89
					&& (resultArray[1] & 0xFF) == 'P'
					&& (resultArray[2] & 0xFF) == 'N'
					&& (resultArray[3] & 0xFF) == 'G') {
				output = Paths.get(filePath + ".png");
			} else {
				// some kind of data file
				output = Paths.get(filePath + ".dat");
			}

			Files.write(output, resultArray, StandardOpenOption.CREATE);
			return resultArray;
		}
		// not encoded, ignore that one
		System.out.println(filePath + " is not encoded");
		return new byte[0];
	}

	/*
	 * Convert the next 4 bytes starting from startIndex into an int
	 * (little-endian data).
	 */
	private static int toInt(byte[] bytes, int startIndex) {
		ByteBuffer buf = ByteBuffer.wrap(bytes, startIndex, 4);
		buf.order(ByteOrder.LITTLE_ENDIAN);
		return buf.getInt();
	}
}
