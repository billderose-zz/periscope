import java.io.*;
import java.net.*;
import java.util.*;

public class AliceInWonderLand {
    private static final String ResourceUrl = "https://s3-us-west-2.amazonaws.com/periscope-public/alice-in-wonderland.txt";

    public static String read() throws Exception {
        final StringBuilder result = new StringBuilder();
        final URL url = new URL(ResourceUrl);
        final Scanner scanner = new Scanner(new InputStreamReader(url.openStream()));
        while (scanner.hasNext()) {
            result.append(scanner.nextLine() + "\n");
        }
        scanner.close();
        return result.toString();
    }

    public static List<String> readTokenized() throws Exception {
        final List<String> result = new LinkedList<>();
        final URL url = new URL(ResourceUrl);
        final Scanner scanner = new Scanner(new InputStreamReader(url.openStream()));
        String token;
        while (scanner.hasNext()) {
            token = scanner.next()
                    .replaceAll("[-.]", " ")
                    .replaceAll("[\\p{Punct}\n\t]", "")
                    .toLowerCase();
            Arrays.asList(token.split(" ")).forEach(t -> result.add(t));
        }
        scanner.close();
        return result;
    }

    public static Map<String, Integer> bagOfWords() throws Exception {
        final Map<String, Integer> bag = new HashMap<>();
        readTokenized().forEach(s -> {
            Integer count = bag.get(s);
            bag.put(s, count == null ? 1 : count + 1);
        });
        return bag;
    }

    public static Map<String, Integer> sortedBagOfWords(Map<String, Integer> bagOfWords) throws Exception {
        final Comparator<String> comparator = (final String o1, final String o2) -> {
            final Integer countOne = bagOfWords.get(o1);
            final Integer countTwo = bagOfWords.get(o2);
            if (countOne.equals(countTwo)) {
                return o1.compareTo(o2);
            }
            return countTwo.compareTo(countOne);
        };
        final Map<String, Integer> sortedBag = new TreeMap<>(comparator);
        sortedBag.putAll(bagOfWords);
        return sortedBag;
    }

    public static void main(String[] args) {
        try {
            sortedBagOfWords(bagOfWords()).forEach((k, v) -> {
                final String bar = new String(new char[v]).replace("\0", "*");
                System.out.printf("%s: %s\n", k, bar);
            });
        } catch (Exception e) {
            System.out.println(e);
        }
    }
}
