<?php
class ContentsHasher {
	private const ALGO = 'md4';

	public static function getFileContentsHash( $content ) {

		return hash( self::ALGO, $content );
	}
}
?>