import org.junit.jupiter.api.Test;

import java.io.File;
import java.net.URL;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class TopHitsTest {

    @Test
    public void verifyNumberOfResultsArgumentParsing() {
        int validResults = TopHits.getNumberOfResultsArgument("-r", "5");
        assertEquals(5, validResults);

        assertThrows(IllegalArgumentException.class, () -> TopHits.getNumberOfResultsArgument("-r"));
        assertThrows(IllegalArgumentException.class, () -> TopHits.getNumberOfResultsArgument("-b", "5"));
        assertThrows(IllegalArgumentException.class, () -> TopHits.getNumberOfResultsArgument("-r", "foo"));
    }

    @Test
    public void verifyFileCheckExists() {
        // To get the absolute path, get the URL first then the full file-name.
        final URL testDataFile = getClass().getResource("test-data-correct.txt");
        assertTrue(TopHits.fileExistsAndIsReadable(testDataFile.getFile()));
    }

    @Test
    public void verifyFileCheckNonExistingFile() {
        assertFalse(TopHits.fileExistsAndIsReadable("some-invalid-file.txt"));
    }

    @Test
    public void testProcessLines() {
        assertEquals(new TopHits.ParsedLine("http://test", 5), TopHits.parseLine("http://test 5"));
        assertEquals(new TopHits.ParsedLine("http://test", 5), TopHits.parseLine("  http://test   5  "));
        assertEquals(new TopHits.ParsedLine("http://test", 5), TopHits.parseLine("http://test 5 IGNORED"));

        assertNull(TopHits.parseLine(null));
        assertNull(TopHits.parseLine("http://test INVALID"));
        assertNull(TopHits.parseLine("http://test"));
    }

    @Test
    public void testSortingCorrectInput() {
        final List<String> topResults = TopHits.getTopResults(new File(getClass().getResource("test-data-correct.txt").getFile()), 5);
        assertEquals("http://api.tech.com/item/3", topResults.get(0));
        assertEquals("http://api.tech.com/item/38", topResults.get(4));
    }

    @Test
    public void testSortingWithInvalidInput() {
        final List<String> topResults = TopHits.getTopResults(new File(getClass().getResource("test-data-invalid-input.txt").getFile()), 5);
        assertEquals("http://api.tech.com/item/7", topResults.get(0));
        assertEquals("http://api.tech.com/item/9", topResults.get(4));
    }

    @Test
    public void testSortingWithShortInput() {
        final List<String> topResults = TopHits.getTopResults(new File(getClass().getResource("test-data-small-file.txt").getFile()), 10);
        assertEquals("http://api.tech.com/item/3", topResults.get(0));
        assertEquals(5, topResults.size());
    }

    @Test
    public void testSortingWithDuplicates() {
        final List<String> topResults = TopHits.getTopResults(new File(getClass().getResource("test-data-duplicates.txt").getFile()), 10);
        assertEquals(10, topResults.size());
        assertEquals("http://api.tech.com/item/3", topResults.get(0));
        // last 6 hits should all be the same value
        for (int i = 4; i < topResults.size(); i++) {
            assertEquals("http://api.tech.com/item/5", topResults.get(i));
        }
    }

}
