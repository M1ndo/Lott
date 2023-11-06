from bs4 import BeautifulSoup
import sys

def parse_html_file(file_path):
    with open(file_path, 'r') as file:
        html_content = file.read()

    soup = BeautifulSoup(html_content, 'html.parser')

    sup_divs = soup.find_all('div', class_='SUP')
    results = {}
    for sup_div in sup_divs:
      p_elements = sup_div.find_all('p')
      if len(p_elements) >= 2:
          date_element = p_elements[1]  # Get the second <p> element
          sorted_element = p_elements[2]  # Get the second <p> element
          sorteded = sorted_element.get_text(strip=True)  # Get the second <p> element
          date = date_element.get_text(strip=True)
          dateAndWhen = date + sorteded
          ul_element = sup_div.find('ul', attrs={'aria-label': 'Números del Super 11'})
          if ul_element:
              numbers = [li.get_text(strip=True) for li in ul_element.find_all('li')]
              results[dateAndWhen] = numbers

    return results

ul_values = parse_html_file(sys.argv[1])
print(ul_values)
