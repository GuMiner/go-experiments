using Newtonsoft.Json;
using System.Collections.Generic;
using System.Configuration;
using System.IO;
using System.Linq;

namespace ReceiptEditor
{
    public class ImageCategory
    {
        [JsonProperty]
        public string Name { get; set; }

        [JsonProperty]
        public string Folder { get; set; }

        public override string ToString()
        {
            return this.Name;
        }

        public static List<ImageCategory> LoadImageCategories()
        {
            return JsonConvert.DeserializeObject<List<ImageCategory>>(File.ReadAllText(GetImageCategoryFilePath())).OrderBy(item => item.Name).ToList();
        }
        
        public static void SaveImageCategories(List<ImageCategory> categories)
        {
            File.WriteAllText(GetImageCategoryFilePath(), JsonConvert.SerializeObject(categories, Formatting.Indented));
        }

        private static string GetImageCategoryFilePath() => ConfigurationManager.AppSettings["CategoriesFile"];
    }
}
