import asyncio
import time
import uuid
import psycopg2
from psycopg2.extras import execute_values
from datetime import datetime

import os
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler, FileCreatedEvent

from utils.to_xml_converter import CSVtoXMLConverter

def get_csv_files_in_input_folder():
    return [os.path.join(dp, f) for dp, dn, filenames in os.walk(CSV_INPUT_PATH) for f in filenames if
            os.path.splitext(f)[1] == '.csv']

def generate_unique_file_name(directory):
    return f"{directory}/{str(uuid.uuid4())}.xml"

def convert_csv_to_xml(in_path, out_path):
    converter = CSVtoXMLConverter(in_path)
    file = open(out_path, "w")
    file.write(converter.to_xml_str())

# Function to connect to db
def connect_to_database():
    try:
        connection = psycopg2.connect(user="is", 
                                    password="is", 
                                    host="db-xml", 
                                    port="5432", 
                                    database="is")
        cursor = connection.cursor()
        return connection, cursor

    except (Exception, psycopg2.Error) as error:
        print("Error connecting to the database:", error)
        return None, None
            
# Function to read XML file
def read_xml_file(file_path):
    with open(file_path, 'r', encoding='utf-8') as file:
        xml_content = file.read()
    return xml_content
        
# Function to get file size
def get_file_size(file_path):
    try:
        return os.path.getsize(file_path)
    except Exception as e:
        print(f"Error getting file size: {e}")
        return None

class CSVHandler(FileSystemEventHandler):
    def __init__(self, input_path, output_path):
        self._output_path = output_path
        self._input_path = input_path

        # generate file creation events for existing files
        for file in [os.path.join(dp, f) for dp, dn, filenames in os.walk(input_path) for f in filenames]:
            event = FileCreatedEvent(os.path.join(CSV_INPUT_PATH, file))
            event.event_type = "created"
            self.dispatch(event)

    async def convert_csv(self, csv_path):
        # here we avoid converting the same file again
        # !TODO: check converted files in the database
        if csv_path in await self.get_converted_files():
            return

        print(f"new file to convert: '{csv_path}'")

        # we generate a unique file name for the XML file
        xml_path = generate_unique_file_name(self._output_path)

        # we do the conversion
        convert_csv_to_xml(csv_path, xml_path)
        print(f"new xml file generated: '{xml_path}'")
        
        # !TODO: once the conversion is done, we should updated the converted_documents tables
        def insert_xml_into_converted_documents(csv_path, xml_path):
            connection, cursor = connect_to_database()

            if connection and cursor:
                try:
                    # Get the file size
                    file_size = get_file_size(xml_path)
                    
                    # Insert XML into the database
                    cursor.execute("INSERT INTO public.converted_documents (src, file_size, dst) VALUES (%s, %s, %s) RETURNING id;",
                                (csv_path, file_size, xml_path))
                    
                    # Commit the transaction
                    connection.commit()

                    print("File added to the converted_documents tables.")

                except (Exception, psycopg2.Error) as error:
                    print("Failed to insert data", error)

                finally:
                    if connection:
                        cursor.close()
                        connection.close()

        # !TODO: we should store the XML document into the imported_documents table
        def insert_xml_into_imported_documents(xml_path):
            connection, cursor = connect_to_database()

            if connection and cursor:
                try:
                    # Read XML file
                    xml_content = read_xml_file(xml_path)
            
                    # Insert XML into the database
                    cursor.execute("INSERT INTO public.imported_documents (file_name, xml) VALUES (%s, %s) RETURNING id;",
                                ('jobdescriptions.xml', xml_content))
                    
                    # Commit the transaction
                    connection.commit()

                    print("XML file added to the imported_documents table.")

                except (Exception, psycopg2.Error) as error:
                    print("Failed to insert data", error)

                finally:
                    if connection:
                        cursor.close()
                        connection.close()
             
        # Call the insert functions
        insert_xml_into_converted_documents(csv_path, xml_path)                
        insert_xml_into_imported_documents(xml_path)
        
        #Soft delete
        """ def soft_delete(xml_path):
            connection, cursor = connect_to_database()
            deleted_on_value = datetime.now()

            try:
                cursor.execute("UPDATE imported_documents SET deleted_on = %s  WHERE file_name = %s", (deleted_on_value, 'jobdescriptions.xml'))
                connection.commit()
                return True
            except Exception as e:
                print(f"Error: {e}")
                connection.rollback()
            finally:
                cursor.close()
                connection.close() """

    async def get_converted_files(self):
        # !TODO: you should retrieve from the database the files that were already converted before
        connection, cursor = connect_to_database()

        try:
            if connection and cursor:
                # Query the database to retrieve files that have already been converted
                cursor.execute("SELECT src FROM public.converted_documents")
                converted_files = [record[0] for record in cursor.fetchall()]

                print("Converted Files:", converted_files)
                
                return converted_files
                
        except (Exception, psycopg2.Error) as error:
            print("Error fetching converted files from the database:", error)

        finally:
            if connection:
                cursor.close()
                connection.close()

        return []

    def on_created(self, event):
        if not event.is_directory and event.src_path.endswith(".csv"):
            asyncio.run(self.convert_csv(event.src_path))


if __name__ == "__main__":

    CSV_INPUT_PATH = "/csv"
    XML_OUTPUT_PATH = "/xml"

    # create the file observer
    observer = Observer()
    observer.schedule(
        CSVHandler(CSV_INPUT_PATH, XML_OUTPUT_PATH),
        path=CSV_INPUT_PATH,
        recursive=True)
    observer.start()

    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        observer.stop()
        observer.join()
