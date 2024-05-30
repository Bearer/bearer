<?php
class ContentsHasher {
	private const ALGO = 'md4';
	private const Type /* comment */ FOO = 'md4';
	private const /** @noinspection PhpUnusedPrivateFieldInspection */ FLD_FLAGS = 4;

	public static function getFileContentsHash( $content ) {

		return hash( self::ALGO, $content );
	}
}
?>