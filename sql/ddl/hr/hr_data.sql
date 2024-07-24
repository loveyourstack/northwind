
INSERT INTO hr.employee (id, first_name, last_name, job_title, title, date_of_birth, hire_date, address, city, state, postal_code, country_fk, home_phone, notes, name) 
  VALUES (-1, '', '', '', '', '1970-01-01', '2010-01-01', '', '', '', '', -1, '', '', 'None');

INSERT INTO hr.employee (first_name, last_name, job_title, title, date_of_birth, hire_date, address, city, state, postal_code, country_fk, home_phone, notes, name) VALUES
  ('Nancy','Davolio','Sales Representative','Ms.','1971-12-08','2015-05-01','507 - 20th Ave. E. Apt. 2A','Seattle','WA','98122',(SELECT id FROM common.country WHERE iso2 = 'US'),'(206) 555-9857','Education includes a BA in psychology from Colorado State University in 1970.  She also completed The Art of the Cold Call.  Nancy is a member of Toastmasters International.','N Davolio'),
  ('Andrew','Fuller','Vice President, Sales','Dr.','1975-02-19','2015-08-14','908 W. Capital Way','Tacoma','WA','98401',(SELECT id FROM common.country WHERE iso2 = 'US'),'(206) 555-9482','Andrew received his BTS commercial in 1974 and a Ph.D. in international marketing from the University of Dallas in 1981.  He is fluent in French and Italian and reads German.  He joined the company as a sales representative, was promoted to sales manager in January 1992 and to vice president of sales in March 1993.  Andrew is a member of the Sales Management Roundtable, the Seattle Chamber of Commerce, and the Pacific Rim Importers Association.','A Fuller'),
  ('Janet','Leverling','Sales Representative','Ms.','1986-08-30','2015-04-01','722 Moss Bay Blvd.','Kirkland','WA','98033',(SELECT id FROM common.country WHERE iso2 = 'US'),'(206) 555-3412','Janet has a BS degree in chemistry from Boston College (1984).  She has also completed a certificate program in food retailing management.  Janet was hired as a sales associate in 1991 and promoted to sales representative in February 1992.','J Leverling'),
  ('Margaret','Peacock','Sales Representative','Mrs.','1960-09-19','2016-05-03','4110 Old Redmond Rd.','Redmond','WA','98052',(SELECT id FROM common.country WHERE iso2 = 'US'),'(206) 555-8122','Margaret holds a BA in English literature from Concordia College (1958) and an MA from the American Institute of Culinary Arts (1966).  She was assigned to the London office temporarily from July through November 1992.','M Peacock'),
  ('Steven','Buchanan','Sales Manager','Mr.','1978-03-04','2016-10-17','14 Garrett Hill','London','','SW1 8JR',(SELECT id FROM common.country WHERE iso2 = 'UK'),'(71) 555-4848','Steven Buchanan graduated from St. Andrews University, Scotland, with a BSC degree in 1976.  Upon joining the company as a sales representative in 1992, he spent 6 months in an orientation program at the Seattle office and then returned to his permanent post in London.  He was promoted to sales manager in March 1993.  Mr. Buchanan has completed the courses Successful Telemarketing and International Sales Management.  He is fluent in French.','S Buchanan'),
  ('Michael','Suyama','Sales Representative','Mr.','1986-07-02','2016-10-17','Coventry House Miner Rd.','London','','EC2 7JR',(SELECT id FROM common.country WHERE iso2 = 'UK'),'(71) 555-7773','Michael is a graduate of Sussex University (MA, economics, 1983) and the University of California at Los Angeles (MBA, marketing, 1986).  He has also taken the courses Multi-Cultural Selling and Time Management for the Sales Professional.  He is fluent in Japanese and can read and write French, Portuguese, and Spanish.','M Suyama'),
  ('Robert','King','Sales Representative','Mr.','1983-05-29','2017-01-02','Edgeham Hollow Winchester Way','London','','RG1 9SP',(SELECT id FROM common.country WHERE iso2 = 'UK'),'(71) 555-5598','Robert King served in the Peace Corps and traveled extensively before completing his degree in English at the University of Michigan in 1992, the year he joined the company.  After completing a course entitled Selling in Europe, he was transferred to the London office in March 1993.','R King'),
  ('Laura','Callahan','Inside Sales Coordinator','Ms.','1981-01-09','2017-03-05','4726 - 11th Ave. N.E.','Seattle','WA','98105',(SELECT id FROM common.country WHERE iso2 = 'US'),'(206) 555-1189','Laura received a BA in psychology from the University of Washington.  She has also completed a course in business French.  She reads and writes French.','L Callahan'),
  ('Anne','Dodsworth','Sales Representative','Ms.','1989-01-27','2017-11-15','7 Houndstooth Rd.','London','','WG2 7LT',(SELECT id FROM common.country WHERE iso2 = 'UK'),'(71) 555-4444','Anne has a BA degree in English from St. Lawrence College.  She is fluent in French and German.','A Dodsworth');

UPDATE hr.employee SET reports_to_fk = (SELECT id FROM hr.employee WHERE last_name = 'Buchanan') WHERE last_name IN ('Suyama','King','Dodsworth');
UPDATE hr.employee SET reports_to_fk = (SELECT id FROM hr.employee WHERE last_name = 'Fuller') WHERE reports_to_fk IS NULL;
ALTER TABLE hr.employee ALTER COLUMN reports_to_fk SET NOT NULL;

