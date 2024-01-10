import xml.etree.ElementTree as ET


class Country:

    def __init__(self, country, location, latitude, longitude):
        Country.counter += 1
        self._id = Country.counter
        self._country = country
        self._location = location
        self._latitude = latitude
        self._longitude = longitude

    def to_xml(self):
        el = ET.Element("Country")
        el.set("id", str(self._id))
        el.set("country", self._country)
        el.set("location", self._location)
        el.set("latitude", self._latitude)
        el.set("longitude", self._longitude)
        return el

    def get_id(self):
        return self._id

    def __str__(self):
        return f"longitude:{self._longitude}, latitude:{self._latitude}, location: {self._location}, country: {self._country}, id:{self._id}"


Country.counter = 0
