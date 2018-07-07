using System;
using System.Drawing;
using System.Drawing.Imaging;
using System.IO;
using System.Linq;
using System.Windows.Forms;

namespace ReceiptEditor
{
    public partial class ImageEditForm : Form
    {
        private Action<ImageAttributes> editCallback;
        private Func<Bitmap> getImageCallback;
        private Action<int> rotateAction;
        private Action hideCallback;

        public ImageEditForm(int subImageId, Action<ImageAttributes> editCallback, Func<Bitmap> getImageCallback, Action<int> rotateAction, Action hideCallback)
        {
            InitializeComponent();
            this.Text = $"Categorization - {subImageId}";
            this.categoryBox.DataSource = ReceiptEditor.ImageCategories;

            this.editCallback = editCallback;
            this.getImageCallback = getImageCallback;
            this.rotateAction = rotateAction;
            this.hideCallback = hideCallback;
        }

        private void brightnessTrackbar_Scroll(object sender, EventArgs e)
        {
            UpdateImageFromEdits();
        }

        private void contrastTrackBar_Scroll(object sender, EventArgs e)
        {
            UpdateImageFromEdits();
        }

        private void UpdateImageFromEdits()
        {
            ImageAttributes imageAttributes = new ImageAttributes();

            imageAttributes.SetColorMatrix(CreateColorMatrix(brightnessTrackbar.Value, contrastTrackBar.Value),
                ColorMatrixFlag.Default, ColorAdjustType.Bitmap);

            editCallback(imageAttributes);
        }

        private ColorMatrix CreateColorMatrix(int trackbarBrighness, int trackbarContrast)
        {
            float brightness = (float)(trackbarBrighness - 50) / 50.0f;
            float contrast = 1.0f + (float)(trackbarContrast - 50) / 50.0f;
            return new ColorMatrix(
                new[]{
                    new float[] {contrast,0,0,0,0},
                    new float[] {0, contrast, 0,0,0},
                    new float[] {0,0, contrast, 0,0},
                    new float[] {0,0,0,1,0},
                    new float[] {brightness, brightness, brightness, 0,1}});
        }

        private void saveButton_Click(object sender, EventArgs e)
        {
            string saveFolder = (this.categoryBox.SelectedItem as ImageCategory).Folder;

            int num = 0;
            string fileName = Path.Combine(saveFolder, $"{num}.jpg");
            while (File.Exists(fileName))
            {
                ++num;
                fileName = Path.Combine(saveFolder, $"{num}.jpg");
            }

            if (!Directory.Exists(Path.GetDirectoryName(fileName)))
            {
                Directory.CreateDirectory(Path.GetDirectoryName(fileName));
            }

            Bitmap imageToSave = this.getImageCallback();
            imageToSave.Save(fileName);
            imageToSave.Dispose();

            ++ReceiptEditor.FilesSharded;

            this.Hide();
            this.hideCallback();
        }

        private void addCategoryButton_Click(object sender, EventArgs e)
        {
            ImageCategory category = new ImageCategory()
            {
                Name = this.newCategoryName.Text,
                Folder = this.newCategoryFolder.Text
            };

            // TODO: This only updates the current subimage window.
            ReceiptEditor.ImageCategories.Add(category);
            ReceiptEditor.ImageCategories.Sort((left, right) => left.Name.CompareTo(right.Name));

            categoryBox.DataSource = null;
            categoryBox.DataSource = ReceiptEditor.ImageCategories;
            categoryBox.DisplayMember = ReceiptEditor.ImageCategories.First().ToString();
            ImageCategory.SaveImageCategories(ReceiptEditor.ImageCategories);
        }

        private void rotateButton_Click(object sender, EventArgs e)
        {
            this.rotateAction(1);
        }
    }
}
